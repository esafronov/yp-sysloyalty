package usecase

import (
	"context"

	"github.com/esafronov/yp-sysloyalty/internal/domain"
	"github.com/esafronov/yp-sysloyalty/internal/helpers/token"
)

type registrationUsecase struct {
	cr domain.CustomerRepository
}

func NewRegisterUsecase(cr domain.CustomerRepository) *registrationUsecase {
	return &registrationUsecase{
		cr: cr,
	}
}

// check if login is busy
func (ru *registrationUsecase) LoginExists(ctx context.Context, login string) (exists bool, err error) {
	exists = true
	c, err := ru.cr.GetByLogin(ctx, login)
	if err != nil {
		return
	}
	if c == nil {
		exists = false
		return
	}
	return
}

// creates new user
func (ru *registrationUsecase) CreateUser(ctx context.Context, customer *domain.Customer) error {
	if err := ru.cr.Create(ctx, customer); err != nil {
		return err
	}
	return nil
}

// creates JWT AccessToken
func (ru *registrationUsecase) CreateAccessToken(user *domain.Customer, secret string, expiry int) (accessToken string, err error) {
	return token.CreateAccessToken(user, secret, expiry)
}
