package controller

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/esafronov/yp-sysloyalty/internal/app/config"
	"github.com/esafronov/yp-sysloyalty/internal/domain"
	"github.com/esafronov/yp-sysloyalty/internal/helpers"
	"github.com/esafronov/yp-sysloyalty/internal/usecase"
)

type OrderController struct {
	or     domain.OrderRepository
	params *config.AppParams
}

func NewOrderController(or domain.OrderRepository, params *config.AppParams) *OrderController {
	return &OrderController{
		or:     or,
		params: params,
	}
}

func (c *OrderController) PostOrder(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "text/plain" {
		http.Error(res, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	orderNum := string(body)
	if !helpers.ValidateOnlyDigits(orderNum) {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err := goluhn.Validate(orderNum); err != nil {
		http.Error(res, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}
	uc := usecase.NewOrdersUsecase(c.or)
	customerID := req.Context().Value(domain.CustomerIDKey).(int64)
	if err = uc.CreateNewOrder(req.Context(), customerID, orderNum); err != nil {
		if errors.Is(err, usecase.ErrOrdersDuplicateOthers) {
			http.Error(res, http.StatusText(http.StatusConflict), http.StatusConflict)
		} else if errors.Is(err, usecase.ErrOrdersDuplicateOwn) {
			res.WriteHeader(http.StatusOK)
		} else {
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	res.WriteHeader(http.StatusAccepted)
}

func (c *OrderController) GetOrders(res http.ResponseWriter, req *http.Request) {
	customerID := req.Context().Value(domain.CustomerIDKey).(int64)
	ou := usecase.NewOrdersUsecase(c.or)
	orders, err := ou.GetOrdersByCustomer(req.Context(), customerID)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	marshaledOrders, err := json.Marshal(orders)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(marshaledOrders)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
