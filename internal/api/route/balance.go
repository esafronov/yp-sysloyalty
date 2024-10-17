package route

import (
	"database/sql"

	"github.com/esafronov/yp-sysloyalty/internal/api/middleware"
	"github.com/esafronov/yp-sysloyalty/internal/app/config"
	"github.com/go-chi/chi"
)

func NewBalanceRoute(r chi.Router, db *sql.DB, params *config.AppParams) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.JwtAuthMiddleware(*params.AccessTokenSecret))
		NewBalanceGetRoute(r, db, params)
		NewBalanceWithdrawRoute(r, db, params)
	})
}
