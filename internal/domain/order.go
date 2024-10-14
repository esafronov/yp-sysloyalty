package domain

import (
	"context"
	"time"
)

type OrderStatus string

const OrderStatusRegistred OrderStatus = "REGISTRED"
const OrderStatusInvalid OrderStatus = "INVALID"
const OrderStatusProcessing OrderStatus = "PROCESSING"
const OrderStatusProcessed OrderStatus = "PROCESSED"

type Order struct {
	ID         int64       `json:"-"`
	CustomerID int64       `json:"-"`
	Num        string      `json:"number"`
	Status     OrderStatus `json:"status"`
	Accrual    int64       `json:"accrual"`
	UploadedAt time.Time   `json:"uploaded_at"`
}

type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetByNum(ctx context.Context, num string) (*Order, error)
	GetAll(ctx context.Context) ([]Order, error)
}
