package route

import (
	"database/sql"

	"github.com/esafronov/yp-sysloyalty/internal/api/controller"
	"github.com/esafronov/yp-sysloyalty/internal/app/config"
	"github.com/esafronov/yp-sysloyalty/internal/repository"
	"github.com/go-chi/chi"
)

func NewRegisterRoute(r chi.Router, db *sql.DB, params *config.AppParams) {
	cr := repository.NewCustomerRepository(db)
	c := controller.NewRegisterController(cr, params)
	r.Post("/register", c.Register)
}
