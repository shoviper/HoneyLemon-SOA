package models

type AccountInfo struct {
	ID       string
	ClientID string
	Type     string
	Balance  float64
	Pin      string
}

type CreateAccount struct {
	Type    string  `json:"type"`
	Balance float64 `json:"balance"`
	Pin     string  `json:"pin"`
}

type AccountVerify struct {
	ID  string `json:"id"`
	Pin string `json:"pin"`
}

type AccountBalance struct {
	ID      string  `json:"id"`
	Balance float64 `json:"balance"`
	Type	string  `json:"type"`
}
