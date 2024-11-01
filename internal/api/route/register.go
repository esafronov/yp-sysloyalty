package route

import (
	"github.com/esafronov/yp-sysloyalty/internal/api/controller"
	"github.com/esafronov/yp-sysloyalty/internal/app/config"
	"github.com/esafronov/yp-sysloyalty/internal/domain"
	"github.com/go-chi/chi"
)

func NewRegisterRoute(r chi.Router, cr domain.CustomerRepository, params *config.AppParams) {
	c := controller.NewRegisterController(cr, params)
	r.Post("/register", c.Register)
}
