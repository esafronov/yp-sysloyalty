package usecase

import (
	"context"

	"github.com/esafronov/yp-sysloyalty/internal/domain"
)

type withdrawlsUsecase struct {
	wr domain.WithdrawRepository
}

func NewWithdrawsUsecase(wr domain.WithdrawRepository) *withdrawlsUsecase {
	return &withdrawlsUsecase{
		wr: wr,
	}
}

func (wc *withdrawlsUsecase) GetWithdrawsByCustomer(ctx context.Context, userID int64) (withdraws []*domain.Withdraw, err error) {
	return wc.wr.GetByCustomer(ctx, userID)
}
