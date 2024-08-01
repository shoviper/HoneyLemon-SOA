package models

import "time"

type ClientInfo struct {
	ID        string
	Name      string
	Address   string
	BirthDate time.Time
}
