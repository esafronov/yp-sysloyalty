package repository

import (
	"context"
	"database/sql"

	"github.com/esafronov/yp-sysloyalty/internal/domain"
)

const WithdrawTable string = "withdrawls"

type withdrawRepository struct {
	db    *sql.DB
	table string
}

func NewWithdrawRepository(db *sql.DB) (r *withdrawRepository, err error) {
	r = &withdrawRepository{
		db:    db,
		table: WithdrawTable,
	}
	err = r.createTable()
	return
}

func (r *withdrawRepository) createTable() error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	// roll back if commit will fail
	defer tx.Rollback()
	tx.Exec(`CREATE TABLE IF NOT EXISTS ` +
		r.table +
		`(
			id bigserial NOT NULL,
			order_num character varying(11) NOT NULL,
			sum integer NOT NULL,
			customer_id bigint NOT NULL,
			processed_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
		)`)
	tx.Exec(`CREATE INDEX IF NOT EXISTS customer_id ON ` + r.table + ` (customer_id)`)
	tx.Exec(`CREATE INDEX IF NOT EXISTS processed_at ON ` + r.table + ` (processed_at)`)
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *withdrawRepository) Create(ctx context.Context, withdraw *domain.Withdraw) error {
	var lastInsertId int64
	row := r.db.QueryRowContext(ctx, "INSERT INTO "+r.table+"(customer_id, order_num, sum) VALUES ($1, $2, $3) RETURNING id", withdraw.CustomerID, withdraw.OrderNum, withdraw.Sum)
	if err := row.Scan(&lastInsertId); err != nil {
		return err
	}
	withdraw.ID = lastInsertId
	return nil
}

func (r *withdrawRepository) GetByCustomer(ctx context.Context, customerID int64) (withdrawls []*domain.Withdraw, err error) {
	withdrawls = make([]*domain.Withdraw, 0)
	rows, err := r.db.QueryContext(ctx, "SELECT id, customer_id, order_num, sum, processed_at FROM "+r.table+" WHERE customer_id=$1 ORDER BY processed_at DESC", customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var withdraw domain.Withdraw
		err = rows.Scan(&withdraw.ID, &withdraw.CustomerID, &withdraw.OrderNum, &withdraw.Sum, &withdraw.ProcessedAt)
		if err != nil {
			return
		}
		withdrawls = append(withdrawls, &withdraw)
	}
	err = rows.Err()
	return
}
