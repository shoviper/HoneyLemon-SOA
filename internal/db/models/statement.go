package models

type StatementInfo struct {
	AccountId string  `json:"accountId"`
	Amount    float64 `json:"amount"`
	Type      string  `json:"type"`
}

type StatementRequest struct {
	AccountId string `json:"accountId"`
}