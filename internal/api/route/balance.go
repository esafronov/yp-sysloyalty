package route

import (
	"github.com/esafronov/yp-sysloyalty/internal/api/middleware"
	"github.com/esafronov/yp-sysloyalty/internal/app/config"
	"github.com/esafronov/yp-sysloyalty/internal/domain"
	"github.com/go-chi/chi"
)

func NewBalanceRoute(r chi.Router, cr domain.CustomerRepository, params *config.AppParams) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.JwtAuthMiddleware(*params.AccessTokenSecret))
		NewBalanceGetRoute(r, cr, params)
		NewBalanceWithdrawRoute(r, cr, params)
	})
}
