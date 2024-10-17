package app

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/esafronov/yp-sysloyalty/internal/api/middleware"
	"github.com/esafronov/yp-sysloyalty/internal/api/route"
	"github.com/esafronov/yp-sysloyalty/internal/app/config"
	"github.com/esafronov/yp-sysloyalty/internal/logger"
	"github.com/esafronov/yp-sysloyalty/internal/postgre"
	"github.com/esafronov/yp-sysloyalty/internal/repository"
	"github.com/go-chi/chi"
)

func Run() error {
	params, err := config.GetAppParams()
	if err != nil {
		return err
	}
	fmt.Println("RunAddress: ", *params.RunAddress)
	fmt.Println("DatabaseURI: ", *params.DatabaseURI)

	if err := postgre.Connect(params.DatabaseURI); err != nil {
		return err
	}

	err = logger.Initialize("debug")
	if err != nil {
		return err
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestLogger(logger.Log))
	r.Use(middleware.GzipCompressing)

	customerRepository, err := repository.NewCustomerRepository(postgre.DB)
	if err != nil {
		fmt.Println("cust rep")
		return err
	}
	orderRepository, err := repository.NewOrderRepository(postgre.DB)
	if err != nil {
		fmt.Println("order rep")
		return err
	}
	withdrawRepisitory, err := repository.NewWithdrawRepository(postgre.DB)
	if err != nil {
		fmt.Println("withd rep")
		return err
	}

	r.Route("/api/user", func(r chi.Router) {
		route.NewRegisterRoute(r, customerRepository, params)
		route.NewLoginRoute(r, customerRepository, params)
		route.NewOrdersRoute(r, orderRepository, params)
		route.NewWithdrawlsRoute(r, withdrawRepisitory, params)
		route.NewBalanceRoute(r, customerRepository, params)
	})

	srv := http.Server{
		Addr:    *params.RunAddress,
		Handler: r,
	}
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
		s := <-sigs
		fmt.Println("got signal ", s)
		srv.Close()
	}()

	return srv.ListenAndServe()
}
