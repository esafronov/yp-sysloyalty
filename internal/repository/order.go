package repository

import (
	"context"
	"database/sql"

	"github.com/esafronov/yp-sysloyalty/internal/domain"
)

const OrderTable string = "orders"

type orderRepository struct {
	db    *sql.DB
	table string
}

func NewOrderRepository(db *sql.DB) (r *orderRepository, err error) {
	r = &orderRepository{
		db:    db,
		table: OrderTable,
	}
	err = r.createTable()
	return
}

func (r *orderRepository) createTable() error {
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
			customer_id bigint NOT NULL,
			order_num character varying(11) NOT NULL,
			accrual integer NOT NULL DEFAULT 0,
			uploaded_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
			status character varying NOT NULL,
			CONSTRAINT order_id PRIMARY KEY (id),
			CONSTRAINT order_num UNIQUE (order_num)
		)`)
	tx.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS order_num ON ` + r.table + ` (order_num)`)
	tx.Exec(`CREATE INDEX IF NOT EXISTS customer_id ON ` + r.table + ` (customer_id)`)
	tx.Exec(`CREATE INDEX IF NOT EXISTS uploaded_at ON ` + r.table + ` (uploaded_at)`)
	tx.Exec(`CREATE INDEX IF NOT EXISTS status ON ` + r.table + ` (status)`)
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *orderRepository) Create(ctx context.Context, order *domain.Order) error {
	var lastInsertID int64
	row := r.db.QueryRowContext(ctx, "INSERT INTO "+r.table+"(customer_id, order_num,status) VALUES ($1, $2, $3) RETURNING id", order.CustomerID, order.Num, order.Status)
	if err := row.Scan(&lastInsertID); err != nil {
		return err
	}
	order.ID = lastInsertID
	return nil
}

func (r *orderRepository) GetByNum(ctx context.Context, num string) (order *domain.Order, err error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, customer_id, order_num, accrual, uploaded_at, status FROM "+r.table+" WHERE order_num=$1 LIMIT 1", num)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, nil
	}
	order = &domain.Order{}
	err = rows.Scan(&order.ID, &order.CustomerID, &order.Num, &order.Accrual, &order.UploadedAt, &order.Status)
	if err != nil {
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (r *orderRepository) GetByCustomer(ctx context.Context, customerID int64) (orders []*domain.Order, err error) {
	orders = make([]*domain.Order, 0)
	rows, err := r.db.QueryContext(ctx, "SELECT id, customer_id, order_num, accrual, uploaded_at, status FROM "+r.table+" WHERE customer_id=$1 ORDER BY uploaded_at DESC", customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var order domain.Order
		err = rows.Scan(&order.ID, &order.CustomerID, &order.Num, &order.Accrual, &order.UploadedAt, &order.Status)
		if err != nil {
			return
		}
		orders = append(orders, &order)
	}
	err = rows.Err()
	return
}

func (r *orderRepository) GetNotFinalStatus(ctx context.Context, limit int) (orders []*domain.Order, err error) {
	orders = make([]*domain.Order, 0)
	rows, err := r.db.QueryContext(ctx, "SELECT id, customer_id, order_num, accrual, uploaded_at, status FROM "+r.table+" WHERE status IN ($1,$2) ORDER BY status, uploaded_at LIMIT $3", domain.OrderStatusRegistred, domain.OrderStatusProcessing, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var order domain.Order
		err = rows.Scan(&order.ID, &order.CustomerID, &order.Num, &order.Accrual, &order.UploadedAt, &order.Status)
		if err != nil {
			return
		}
		orders = append(orders, &order)
	}
	err = rows.Err()
	return
}

func (r *orderRepository) UpdateStatus(ctx context.Context, num string, status domain.OrderStatus) (err error) {
	_, err = r.db.ExecContext(ctx, "UPDATE "+r.table+" SET status=$1 WHERE order_num=$2", status, num)
	if err != nil {
		return
	}
	return
}
