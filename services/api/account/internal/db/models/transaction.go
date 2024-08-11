package models

import "time"

type TransactionInfo struct {
	ID         string
	SenderID   string
	ReceiverID string
	Amount     float64
	CreatedAt  time.Time
}
