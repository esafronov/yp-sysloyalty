package controller

import (
	"encoding/json"
	"net/http"

	"github.com/esafronov/yp-sysloyalty/internal/app/config"
	"github.com/esafronov/yp-sysloyalty/internal/domain"
	"github.com/esafronov/yp-sysloyalty/internal/usecase"
	"golang.org/x/crypto/bcrypt"
)

type LoginController struct {
	cr     domain.CustomerRepository
	params *config.AppParams
}

func NewLoginController(cr domain.CustomerRepository, params *config.AppParams) *LoginController {
	return &LoginController{
		cr:     cr,
		params: params,
	}
}

func (c *LoginController) Login(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "application/json" {
		http.Error(res, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}
	var request domain.LoginRequest

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	uc := usecase.NewLoginUsecase(c.cr)
	customer, err := uc.FindUserByLogin(req.Context(), request.Login)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if customer == nil {
		http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(request.Password))
	if err != nil {
		http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	accessToken, err := uc.CreateAccessToken(customer, *c.params.AccessTokenSecret, *c.params.ExpireAccessToken)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	refreshToken, err := uc.CreateRefreshToken(customer, *c.params.RefreshTokenSecret, *c.params.ExpireRefreshToken)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	reponse := &domain.LoginReponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	responseJson, err := json.Marshal(reponse)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(responseJson)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
