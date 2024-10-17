package usecase

import (
	"context"
	"errors"

	"github.com/esafronov/yp-sysloyalty/internal/domain"
)

var (
	ErrWithdrawInsufficientBalance = errors.New("insufficient balance")
)

type withdrawUsecase struct {
	cr domain.CustomerRepository
}

func NewWithdrawUsecase(cr domain.CustomerRepository) *withdrawUsecase {
	return &withdrawUsecase{
		cr: cr,
	}
}

func (wc *withdrawUsecase) Withdraw(ctx context.Context, userID int64, req *domain.WithdrawRequest) error {
	return wc.cr.Withdraw(ctx, userID, req.OrderNum, req.Sum, func(customer *domain.Customer) error {
		if !customer.CanWithdraw(req.Sum) {
			return ErrWithdrawInsufficientBalance
		}
		return nil
	})
}
