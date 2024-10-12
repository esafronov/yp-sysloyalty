package repository

import (
	"context"
	"database/sql"

	"github.com/esafronov/yp-sysloyalty/internal/domain"
)

const CustomerTable string = "customers"

type customerRepository struct {
	db    *sql.DB
	table string
}

func NewCustomerRepository(db *sql.DB) *customerRepository {
	return &customerRepository{
		db:    db,
		table: CustomerTable,
	}
}

func (r *customerRepository) Create(ctx context.Context, user *domain.Customer) error {
	var lastInsertId int64
	row := r.db.QueryRowContext(ctx, "INSERT INTO "+r.table+"(login, password) VALUES ($1, $2) RETURNING id", user.Login, user.Password)
	if err := row.Scan(&lastInsertId); err != nil {
		return err
	}
	user.ID = lastInsertId
	return nil
}

func (r *customerRepository) CreditBalance(ctx context.Context, userID int64, amount int) error {
	return nil
}

func (r *customerRepository) DebitBalance(ctx context.Context, userID int64, amount int) error {
	return nil
}

func (r *customerRepository) GetByLogin(ctx context.Context, login string) (*domain.Customer, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, login, password, balance, withdrawn FROM "+r.table+" WHERE login=$1 LIMIT 1", login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, nil
	}
	var customer domain.Customer
	err = rows.Scan(&customer.ID, &customer.Login, &customer.Password, &customer.Balance, &customer.Withdrawn)
	if err != nil {
		return nil, err
	}
	return &customer, nil
}
