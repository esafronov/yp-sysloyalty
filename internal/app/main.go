package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/esafronov/yp-sysloyalty/internal/api/middleware"
	"github.com/esafronov/yp-sysloyalty/internal/api/route"
	"github.com/esafronov/yp-sysloyalty/internal/app/config"
	"github.com/esafronov/yp-sysloyalty/internal/logger"
	"github.com/esafronov/yp-sysloyalty/internal/postgre"
	"github.com/esafronov/yp-sysloyalty/internal/repository"
	"github.com/esafronov/yp-sysloyalty/internal/routine"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func Run() error {
	params, err := config.GetAppParams()
	if err != nil {
		return err
	}

	logger.Log.Info("params",
		zap.String("RunAddress", *params.RunAddress),
		zap.String("DatabaseURI", *params.DatabaseURI),
		zap.String("AccrualSystemAddress", *params.AccrualSystemAddress),
		zap.String("AccessTokenSecret", *params.AccessTokenSecret),
		zap.Int("ExpireAccessToken", *params.ExpireAccessToken),
		zap.Int("GrabInterval", *params.GrabInterval),
		zap.Int("ProcessRate", *params.ProcessRate),
	)

	if err := postgre.Connect(params.DatabaseURI); err != nil {
		return err
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestLogger(logger.Log))
	r.Use(middleware.GzipCompressing)

	customerRepository, err := repository.NewCustomerRepository(postgre.DB)
	if err != nil {
		return err
	}
	orderRepository, err := repository.NewOrderRepository(postgre.DB)
	if err != nil {
		return err
	}
	withdrawRepisitory, err := repository.NewWithdrawRepository(postgre.DB)
	if err != nil {
		return err
	}

	r.Route("/api/user", func(r chi.Router) {
		route.NewRegisterRoute(r, customerRepository, params)
		route.NewLoginRoute(r, customerRepository, params)
		route.NewOrdersRoute(r, orderRepository, params)
		route.NewWithdrawlsRoute(r, withdrawRepisitory, params)
		route.NewBalanceRoute(r, customerRepository, params)
	})

	ctx, appExit := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	poller := routine.NewPoller(params)
	grabber := routine.NewGrabber(orderRepository, params)

	orderChan := grabber.Run(ctx, poller.RetryAfterChan, &wg)
	updateChan := poller.Run(ctx, orderChan, &wg)

	updater := routine.NewUpdater(orderRepository, customerRepository, params)
	updater.Run(ctx, updateChan, &wg)

	srv := http.Server{
		Addr:    *params.RunAddress,
		Handler: r,
	}
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
		s := <-sigs
		fmt.Println("got signal ", s)
		appExit()
		wg.Wait()
		srv.Close()
	}()

	return srv.ListenAndServe()
}
