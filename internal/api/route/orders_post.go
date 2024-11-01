package route

import (
	"github.com/esafronov/yp-sysloyalty/internal/api/controller"
	"github.com/esafronov/yp-sysloyalty/internal/app/config"
	"github.com/esafronov/yp-sysloyalty/internal/domain"
	"github.com/go-chi/chi"
)

func NewOrdersPostRoute(r chi.Router, or domain.OrderRepository, params *config.AppParams) {
	c := controller.NewOrderController(or, params)
	r.Post("/orders", c.PostOrder)
}
