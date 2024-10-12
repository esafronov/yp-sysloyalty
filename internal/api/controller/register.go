package controller

import (
	"encoding/json"
	"net/http"

	"github.com/esafronov/yp-sysloyalty/internal/app/config"
	"github.com/esafronov/yp-sysloyalty/internal/domain"
	"github.com/esafronov/yp-sysloyalty/internal/logger"
	"github.com/esafronov/yp-sysloyalty/internal/usecases"
	"go.uber.org/zap"
)

type RegisterController struct {
	cr     domain.CustomerRepository
	params *config.AppParams
}

func NewRegisterController(cr domain.CustomerRepository, params *config.AppParams) *RegisterController {
	return &RegisterController{
		cr:     cr,
		params: params,
	}
}

func (c *RegisterController) Register(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "application/json" {
		http.Error(res, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}
	var request domain.RegistrationRequest

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	uc := usecases.NewRegisterUsecase(c.cr)
	exists, err := uc.LoginExists(req.Context(), request.Login)
	if err != nil {
		logger.Log.Error("error ", zap.Error(err))
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(res, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}
	customer := &domain.Customer{
		Login:    request.Login,
		Password: request.Password,
	}
	if err := uc.CreateUser(req.Context(), customer); err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
	reponse := &domain.RegistrationReponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	responseJson, err := json.Marshal(reponse)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	_, err = res.Write(responseJson)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}
