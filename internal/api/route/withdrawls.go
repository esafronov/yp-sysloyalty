package route

import (
	"database/sql"

	"github.com/esafronov/yp-sysloyalty/internal/api/controller"
	"github.com/esafronov/yp-sysloyalty/internal/api/middleware"
	"github.com/esafronov/yp-sysloyalty/internal/app/config"
	"github.com/esafronov/yp-sysloyalty/internal/repository"
	"github.com/go-chi/chi"
)

func NewWithdrawlsRoute(r chi.Router, db *sql.DB, params *config.AppParams) {
	wr := repository.NewWithdrawRepository(db)
	c := controller.NewWithdrawController(wr, params)
	r.Group(func(r chi.Router) {
		r.Use(middleware.JwtAuthMiddleware(*params.AccessTokenSecret))
		r.Get("/withdrawls", c.Withdrawls)
	})
}
