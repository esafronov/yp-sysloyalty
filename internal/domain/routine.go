package domain

import (
	"encoding/json"
)

type GrabberRetryAfterType string

var GrabberRetryAfterKey GrabberRetryAfterType = "retry-after"

type OrderUpdate struct {
	Num     string      `json:"order"`
	Status  OrderStatus `json:"status"`
	Accrual int64       `json:"accrual,omitempty"`
}

func (o *OrderUpdate) UnmarshalJSON(data []byte) (err error) {
	type OrderUpdateAlias OrderUpdate
	aliasValue := &struct {
		*OrderUpdateAlias
		Accrual float64 `json:"accrual,omitempty"`
	}{
		OrderUpdateAlias: (*OrderUpdateAlias)(o),
	}
	if err := json.Unmarshal(data, aliasValue); err != nil {
		return err
	}
	o.Accrual = int64(aliasValue.Accrual * 100)
	return
}
