package route

import (
	"github.com/esafronov/yp-sysloyalty/internal/api/controller"
	"github.com/esafronov/yp-sysloyalty/internal/app/config"
	"github.com/esafronov/yp-sysloyalty/internal/domain"
	"github.com/go-chi/chi"
)

func NewBalanceGetRoute(r chi.Router, cr domain.CustomerRepository, params *config.AppParams) {
	c := controller.NewBalanceController(cr, params)
	r.Get("/balance", c.GetBalance)
}
