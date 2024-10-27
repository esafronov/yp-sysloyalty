package route

import (
	"github.com/esafronov/yp-sysloyalty/internal/api/controller"
	"github.com/esafronov/yp-sysloyalty/internal/api/middleware"
	"github.com/esafronov/yp-sysloyalty/internal/app/config"
	"github.com/esafronov/yp-sysloyalty/internal/domain"
	"github.com/go-chi/chi"
)

func NewWithdrawlsRoute(r chi.Router, wr domain.WithdrawRepository, params *config.AppParams) {
	c := controller.NewWithdrawController(wr, params)
	r.Group(func(r chi.Router) {
		r.Use(middleware.JwtAuthMiddleware(*params.AccessTokenSecret))
		r.Get("/withdrawals", c.Withdrawls)
	})
}
