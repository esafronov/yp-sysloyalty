package domain

import (
	"context"
	"encoding/json"
)

type Customer struct {
	ID        int64  `json:"-"`
	Login     string `json:"-"`
	Password  string `json:"-"`
	Withdrawn int64  `json:"-"`
	Balance   int64  `json:"-"`
}

type CustomerRepository interface {
	Create(ctx context.Context, user *Customer) error
	GetByLogin(ctx context.Context, login string) (*Customer, error)
	GetByID(ctx context.Context, userID int64) (*Customer, error)
	Withdraw(ctx context.Context, userID int64, orderNum string, sum int64, updateFunc func(customer *Customer) error) error
}

type customerContextKey string

const CustomerIDKey customerContextKey = "x-user-id"

func (c *Customer) MarshalJSON() ([]byte, error) {
	type aliasCustomer Customer
	var balance float64 = float64(c.Balance) / 100
	var withdrawn float64 = float64(c.Withdrawn) / 100
	alias := struct {
		aliasCustomer
		Withdrawn float64 `json:"withdrawn"`
		Balance   float64 `json:"current"`
	}{
		aliasCustomer: aliasCustomer(*c),
		Withdrawn:     withdrawn,
		Balance:       balance,
	}
	return json.Marshal(alias)
}

func (c *Customer) CanWithdraw(sum int64) (result bool) {
	if c.Balance >= sum {
		result = true
	}
	return
}
