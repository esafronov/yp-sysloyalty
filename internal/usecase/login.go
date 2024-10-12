package usecase

import (
	"context"

	"github.com/esafronov/yp-sysloyalty/internal/domain"
	"github.com/esafronov/yp-sysloyalty/internal/helpers/token"
)

type loginUsecase struct {
	cr domain.CustomerRepository
}

func NewLoginUsecase(cr domain.CustomerRepository) *loginUsecase {
	return &loginUsecase{
		cr: cr,
	}
}

func (lu *loginUsecase) FindUserByLogin(ctx context.Context, login string) (customer *domain.Customer, err error) {
	return lu.cr.GetByLogin(ctx, login)
}

func (lu *loginUsecase) CreateAccessToken(user *domain.Customer, secret string, expiry int) (accessToken string, err error) {
	return token.CreateAccessToken(user, secret, expiry)
}

func (lu *loginUsecase) CreateRefreshToken(user *domain.Customer, secret string, expiry int) (refreshToken string, err error) {
	return token.CreateRefreshToken(user, secret, expiry)
}
