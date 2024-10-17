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

func NewOrderRepository(db *sql.DB) *orderRepository {
	return &orderRepository{
		db:    db,
		table: OrderTable,
	}
}

func (r *orderRepository) Create(ctx context.Context, order *domain.Order) error {
	var lastInsertId int64
	row := r.db.QueryRowContext(ctx, "INSERT INTO "+r.table+"(customer_id, order_num,status) VALUES ($1, $2, $3) RETURNING id", order.CustomerID, order.Num, order.Status)
	if err := row.Scan(&lastInsertId); err != nil {
		return err
	}
	order.ID = lastInsertId
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
