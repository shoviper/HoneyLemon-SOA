package models

import "time"

type StatementInfo struct {
	ClientName string         `json:"clientName"`
	ClientID   string         `json:"clientID"`
	AccountID  string         `json:"accountID"`
	Start      time.Time      `json:"start"`
	End        time.Time      `json:"end"`
	Activity   []ActivityInfo `json:"activity"`
}

type ActivityInfo struct {
	TxID      string    `json:"txID"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Amount    float64   `json:"amount"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
}
