package models

import "time"

type ClientInfo struct {
	ID        byte   `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	BirthDate time.Time `json:"birth_date"`
}