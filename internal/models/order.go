package models

import (
	"time"
)

type Order struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	CustomerID int64     `json:"customer_id"`
	Amount     float64   `json:"amount"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type OrderReturn struct {
	ID           int64      `gorm:"primaryKey" json:"id"`
	OrderID      int64      `json:"order_id"`
	CustomerID   int64      `json:"customer_id"`
	Reason       string     `json:"reason"`
	RequestedAt  time.Time  `json:"requested_at"`
	Approved     bool       `json:"approved"`
	ProcessedAt  *time.Time `json:"processed_at,omitempty"`
	RefundAmount float64    `json:"refund_amount"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
