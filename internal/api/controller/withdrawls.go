package controller

import (
	"encoding/json"
	"net/http"

	"github.com/esafronov/yp-sysloyalty/internal/app/config"
	"github.com/esafronov/yp-sysloyalty/internal/domain"
)

type WithdrawController struct {
	wr     domain.WithdrawRepository
	params *config.AppParams
}

func NewWithdrawController(wr domain.WithdrawRepository, params *config.AppParams) *WithdrawController {
	return &WithdrawController{
		wr:     wr,
		params: params,
	}
}

func (c *WithdrawController) Withdrawls(res http.ResponseWriter, req *http.Request) {
	customerID := req.Context().Value(domain.CustomerIDKey).(int64)
	withdrawls, err := c.wr.GetByCustomer(req.Context(), customerID)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	marshaledWithdrawls, err := json.Marshal(withdrawls)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(marshaledWithdrawls)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
