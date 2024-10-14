package route

import (
	"database/sql"

	"github.com/esafronov/yp-sysloyalty/internal/api/controller"
	"github.com/esafronov/yp-sysloyalty/internal/app/config"
	"github.com/esafronov/yp-sysloyalty/internal/repository"
	"github.com/go-chi/chi"
)

func NewOrdersPostRoute(r chi.Router, db *sql.DB, params *config.AppParams) {
	cr := repository.NewOrderRepository(db)
	c := controller.NewOrderController(cr, params)
	r.Post("/orders", c.PostOrder)
}
