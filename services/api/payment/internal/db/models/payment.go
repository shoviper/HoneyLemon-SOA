package models

import "time"

type PaymentInfo struct {
	ID        string
	AccountID string
	RefCode   string
	Amount    float64
	CreatedAt time.Time
}
