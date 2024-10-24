package usecase

import (
	"context"
	"errors"

	"github.com/esafronov/yp-sysloyalty/internal/domain"
)

var (
	ErrOrdersHasFinalStatus = errors.New("order has final status")
)

type ordersUpdateUsecase struct {
	or domain.OrderRepository
	cr domain.CustomerRepository
}

func NewOrdersUpdateUsecase(or domain.OrderRepository, cr domain.CustomerRepository) *ordersUpdateUsecase {
	return &ordersUpdateUsecase{
		or: or,
		cr: cr,
	}
}

func (ou *ordersUpdateUsecase) Update(ctx context.Context, update *domain.OrderUpdate) error {
	o, err := ou.or.GetByNum(ctx, update.Num)
	if err != nil {
		return err
	}
	if o.HasFinalStatus() {
		return ErrOrdersHasFinalStatus
	}
	if update.Accrual > 0 {
		return ou.cr.Accrual(ctx, o.CustomerID, update.Num, update.Accrual)
	}
	return ou.or.UpdateStatus(ctx, update.Num, update.Status)
}
