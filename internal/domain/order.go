package domain

import (
	"context"
	"encoding/json"
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
	Accrual    int64       `json:"accrual,omitempty"`
	UploadedAt time.Time   `json:"uploaded_at"`
}

func (o *Order) HasFinalStatus() bool {
	if o.Status == OrderStatusProcessed || o.Status == OrderStatusInvalid {
		return true
	}
	return false
}

func (o *Order) MarshalJSON() ([]byte, error) {
	formattedDate := o.UploadedAt.Format(time.RFC3339)
	type aliasOrder Order
	alias := struct {
		aliasOrder
		UploadedAt string `json:"uploaded_at"`
	}{
		aliasOrder: aliasOrder(*o),
		UploadedAt: formattedDate,
	}
	return json.Marshal(alias)
}

type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetByNum(ctx context.Context, num string) (*Order, error)
	GetByCustomer(ctx context.Context, customerID int64) ([]*Order, error)
	GetNotFinalStatus(ctx context.Context) ([]*Order, error)
	UpdateStatus(ctx context.Context, num string, status OrderStatus) error
}
