package domain

type GrabberRetryAfterType string

var GrabberRetryAfterKey GrabberRetryAfterType = "retry-after"

type OrderUpdate struct {
	Num     string      `json:"order"`
	Status  OrderStatus `json:"status"`
	Accrual int64       `json:"accrual,omitempty"`
}
