package usecase

import (
	"context"

	"github.com/esafronov/yp-sysloyalty/internal/domain"
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
	orders, err = ou.or.GetByCustomer(ctx, customerID)
	if err != nil {
		return
	}
	return
}
