package usecase

import (
	"context"
	"errors"

	"github.com/esafronov/yp-sysloyalty/internal/domain"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrOrdersDuplicateOwn    = errors.New("order duplicate your own")
	ErrOrdersDuplicateOthers = errors.New("order duplicate others")
)

type ordersUsecase struct {
	or domain.OrderRepository
}

func NewOrdersUsecase(or domain.OrderRepository) *ordersUsecase {
	return &ordersUsecase{
		or: or,
	}
}

func (ou *ordersUsecase) GetOrdersByCustomer(ctx context.Context, customerID int64) (orders []*domain.Order, err error) {
	return ou.or.GetByCustomer(ctx, customerID)
}

func (ou *ordersUsecase) CreateNewOrder(ctx context.Context, customerID int64, orderNum string) (err error) {
repeatOrderLoad:
	order, err := ou.or.GetByNum(ctx, orderNum)
	if err != nil {
		return
	}
	if order != nil {
		if order.CustomerID == customerID {
			err = ErrOrdersDuplicateOwn
		} else {
			err = ErrOrdersDuplicateOthers
		}
		return
	}
	order = &domain.Order{
		CustomerID: customerID,
		Num:        orderNum,
		Status:     domain.OrderStatusRegistred,
	}
	if err = ou.or.Create(ctx, order); err != nil {
		if data, ok := err.(*pgconn.PgError); ok && data.Code == pgerrcode.UniqueViolation {
			goto repeatOrderLoad
		}
	}
	return
}
