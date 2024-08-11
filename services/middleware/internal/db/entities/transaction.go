package entities

import (
	cryptoRand "crypto/rand"
	"math/big"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Transaction struct {
	ID         string    `gorm:"primaryKey"`
	SenderID   string    `gorm:"not null"`
	Sender     Account   `gorm:"foreignKey:SenderID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ReceiverID string    `gorm:"not null"`
	Receiver   Account   `gorm:"foreignKey:ReceiverID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Amount     float64   `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt
}

func (Transaction) TableName() string {
	return "transactions"
}

func (tr *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	t := time.Now().UTC()
	entropy := ulid.Monotonic(cryptoRand.Reader, 0)
	ulid := ulid.MustNew(ulid.Timestamp(t), entropy)

	// Convert the ULID to a big.Int
	bigInt := new(big.Int).SetBytes(ulid[:])

	// Get the last 12 digits
	modulus := new(big.Int).SetInt64(1000000000000)
	result := new(big.Int).Mod(bigInt, modulus)

	// Convert the result to a string
	id := result.String()

	// Ensure the result is 12 digits by padding with leading zeros if necessary
	for len(id) < 10 {
		id = "0" + id
	}

	tr.ID = id
	tr.CreatedAt = time.Now()
	tr.UpdatedAt = time.Now()
	return
}

func (tr *Transaction) BeforeUpdate(tx *gorm.DB) (err error) {
	tr.UpdatedAt = time.Now()
	return
}
