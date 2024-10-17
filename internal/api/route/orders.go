package route

import (
	"github.com/esafronov/yp-sysloyalty/internal/api/middleware"
	"github.com/esafronov/yp-sysloyalty/internal/app/config"
	"github.com/esafronov/yp-sysloyalty/internal/domain"
	"github.com/go-chi/chi"
)

func NewOrdersRoute(r chi.Router, or domain.OrderRepository, params *config.AppParams) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.JwtAuthMiddleware(*params.AccessTokenSecret))
		NewOrdersPostRoute(r, or, params)
		NewOrdersGetRoute(r, or, params)
	})
}
