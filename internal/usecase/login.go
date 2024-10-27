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

// find user by login
func (lu *loginUsecase) FindUserByLogin(ctx context.Context, login string) (customer *domain.Customer, err error) {
	return lu.cr.GetByLogin(ctx, login)
}

// create JWT AccessToken
func (lu *loginUsecase) CreateAccessToken(user *domain.Customer, secret string, expiry int) (accessToken string, err error) {
	return token.CreateAccessToken(user, secret, expiry)
}
