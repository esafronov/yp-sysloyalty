package controller

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/esafronov/yp-sysloyalty/internal/app/config"
	"github.com/esafronov/yp-sysloyalty/internal/domain"
	"github.com/esafronov/yp-sysloyalty/internal/usecase"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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
	if err := goluhn.Validate(orderNum); err != nil {
		http.Error(res, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}
repeatOrderLoad:
	order, err := c.or.GetByNum(req.Context(), orderNum)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	cid := req.Context().Value(domain.CustomerIDKey).(int64)
	if order != nil {
		if order.CustomerID == cid {
			res.WriteHeader(http.StatusOK)
		} else {
			http.Error(res, http.StatusText(http.StatusConflict), http.StatusConflict)
		}
		return
	}
	order = &domain.Order{
		CustomerID: cid,
		Num:        orderNum,
		Status:     domain.OrderStatusRegistred,
	}
	if err = c.or.Create(req.Context(), order); err != nil {
		if data, ok := err.(*pgconn.PgError); ok && data.Code == pgerrcode.UniqueViolation {
			goto repeatOrderLoad
		}
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusAccepted)
}

func (c *OrderController) GetOrders(res http.ResponseWriter, req *http.Request) {
	cid := req.Context().Value(domain.CustomerIDKey).(int64)
	ou := usecase.NewOrdersUsecase(c.or)
	orders, err := ou.GetOrdersByCustomer(req.Context(), cid)
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
