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

func NewWithdrawRepository(db *sql.DB) *withdrawRepository {
	return &withdrawRepository{
		db:    db,
		table: WithdrawTable,
	}
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
