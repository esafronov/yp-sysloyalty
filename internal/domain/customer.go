package domain

import "context"

type Customer struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Password  string `json:"password"`
	Withdrawn int64  `json:"withdrown"`
	Balance   int64  `json:"withdrawn"`
}

type CustomerRepository interface {
	Create(ctx context.Context, user *Customer) error
	CreditBalance(ctx context.Context, userID int64, credit int) error
	DebitBalance(ctx context.Context, userID int64, debit int) error
	GetByLogin(ctx context.Context, login string) (*Customer, error)
}

type customerContextKey string

const CustomerIDKey customerContextKey = "x-user-id"
