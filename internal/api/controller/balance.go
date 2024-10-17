package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/esafronov/yp-sysloyalty/internal/app/config"
	"github.com/esafronov/yp-sysloyalty/internal/domain"
	"github.com/esafronov/yp-sysloyalty/internal/logger"
	"github.com/esafronov/yp-sysloyalty/internal/usecase"
	"go.uber.org/zap"
)

type BalanceController struct {
	cr     domain.CustomerRepository
	params *config.AppParams
}

func NewBalanceController(cr domain.CustomerRepository, params *config.AppParams) *BalanceController {
	return &BalanceController{
		cr:     cr,
		params: params,
	}
}

func (c *BalanceController) Withdraw(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "application/json" {
		http.Error(res, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}
	var request domain.WithdrawRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if request.OrderNum == "" || request.Sum == 0 {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err := goluhn.Validate(request.OrderNum); err != nil {
		http.Error(res, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}
	customerID := req.Context().Value(domain.CustomerIDKey).(int64)
	uc := usecase.NewWithdrawUsecase(c.cr)
	err := uc.Withdraw(req.Context(), customerID, &request)
	if err != nil {
		if errors.Is(err, usecase.ErrWithdrawInsufficientBalance) {
			http.Error(res, err.Error(), http.StatusPaymentRequired)
		} else {
			logger.Log.Error("withdraw :", zap.Error(err))
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
}

func (c *BalanceController) GetBalance(res http.ResponseWriter, req *http.Request) {
	customerID := req.Context().Value(domain.CustomerIDKey).(int64)
	customer, err := c.cr.GetByID(req.Context(), customerID)
	if err != nil || customer == nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	marshaledCustomer, err := json.Marshal(customer)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(marshaledCustomer)
}
