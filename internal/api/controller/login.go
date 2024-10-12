package controller

import (
	"net/http"

	"github.com/esafronov/yp-sysloyalty/internal/domain"
)

type LoginController struct {
	cr domain.CustomerRepository
}

func NewLoginController(cr domain.CustomerRepository) *LoginController {
	return &LoginController{
		cr: cr,
	}
}

func (c *LoginController) Login(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	res.WriteHeader(http.StatusOK)
}
