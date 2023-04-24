package entity

import (
	"time"
)

type Order struct {
	Number     string    `json:"number" db:"order_number"`
	Status     string    `json:"status" db:"status"`
	Accrual    float64   `json:"accrual,omitempty" db:"accrual"`
	UploadedAt time.Time `json:"uploaded_at" db:"uploaded_at"`
	UserID     int       `json:"-" db:"user_id"`
}

type OrderAdapter struct {
	Number  string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}
