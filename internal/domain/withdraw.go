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
	sum := float64(w.Sum) / 100
	type aliasWithdraw Withdraw
	alias := struct {
		aliasWithdraw
		Sum         float64 `json:"sum"`
		ProcessedAt string  `json:"processed_at"`
	}{
		aliasWithdraw: aliasWithdraw(*w),
		ProcessedAt:   formattedDate,
		Sum:           sum,
	}
	return json.Marshal(alias)
}

type WithdrawRequest struct {
	OrderNum string `json:"order"`
	Sum      int64  `json:"sum"`
}

func (wr *WithdrawRequest) UnmarshalJSON(data []byte) (err error) {
	type WithdrawRequestAlias WithdrawRequest
	aliasValue := &struct {
		*WithdrawRequestAlias
		Sum float64 `json:"sum"`
	}{
		WithdrawRequestAlias: (*WithdrawRequestAlias)(wr),
	}
	if err := json.Unmarshal(data, aliasValue); err != nil {
		return err
	}
	wr.Sum = int64(aliasValue.Sum * 100)
	return
}

type WithdrawRepository interface {
	Create(ctx context.Context, withdraw *Withdraw) error
	GetByCustomer(ctx context.Context, customerID int64) ([]*Withdraw, error)
}
