package domain

import (
	"context"
	"encoding/json"
	"time"
)

type Withdraw struct {
	ID          int64     `json:"-"`
	CustomerID  int64     `json:"-"`
	OrderNum    string    `json:"order"`
	Sum         int64     `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

func (w *Withdraw) MarshalJSON() ([]byte, error) {
	formattedDate := w.ProcessedAt.Format(time.RFC3339)
	type aliasWithdraw Withdraw
	alias := struct {
		aliasWithdraw
		ProcessedAt string `json:"processed_at"`
	}{
		aliasWithdraw: aliasWithdraw(*w),
		ProcessedAt:   formattedDate,
	}
	return json.Marshal(alias)
}

type WithdrawRequest struct {
	OrderNum string `json:"order"`
	Sum      int64  `json:"sum"`
}

type WithdrawRepository interface {
	Create(ctx context.Context, withdraw *Withdraw) error
	GetByCustomer(ctx context.Context, customerID int64) ([]*Withdraw, error)
}
