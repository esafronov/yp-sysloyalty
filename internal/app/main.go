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
	//logger middleware
	r.Use(middleware.RequestLogger(logger.Log))
	//compressing middleware
	r.Use(middleware.GzipCompressing)

	//init customer repository
	customerRepository, err := repository.NewCustomerRepository(postgre.DB)
	if err != nil {
		return err
	}

	//init order repository
	orderRepository, err := repository.NewOrderRepository(postgre.DB)
	if err != nil {
		return err
	}

	//init withdraw repository
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

	var wg sync.WaitGroup //waitgroup for 3 main app routines: grabber, poller, updater

	//initialize poller and grabber routines
	poller := routine.NewPoller(params)
	grabber := routine.NewGrabber(orderRepository, params)

	//run grabber and poller routines
	orderChan := grabber.Run(ctx, poller.RetryAfterChan, &wg)
	updateChan := poller.Run(ctx, orderChan, &wg)

	//initialize and run updater routine
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
		appExit() //send app exit signal
		wg.Wait() //wait until all routines are stoped
		srv.Close()
	}()

	return srv.ListenAndServe()
}
