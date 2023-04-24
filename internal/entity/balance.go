package entity

import (
	"time"
)

type Withdrawal struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type Balance struct {
	Current  float64 `json:"current" db:"balance"`
	Withdraw float64 `json:"withdrawn" db:"withdrawal"`
}

type Withdraw struct {
	Order       string    `json:"order" db:"order_number"`
	Sum         float64   `json:"sum" db:"sum_number"`
	ProcessedAt time.Time `json:"processed_at" db:"updated_at"`
}
