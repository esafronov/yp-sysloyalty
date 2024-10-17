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

	r.Route("/api/user", func(r chi.Router) {
		route.NewRegisterRoute(r, postgre.DB, params)
		route.NewLoginRoute(r, postgre.DB, params)
		route.NewOrdersRoute(r, postgre.DB, params)
		route.NewWithdrawlsRoute(r, postgre.DB, params)
		route.NewBalanceRoute(r, postgre.DB, params)
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
