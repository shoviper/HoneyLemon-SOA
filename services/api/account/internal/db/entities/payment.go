package entities

import (
	cryptoRand "crypto/rand"
	"math/big"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Payment struct {
	ID        string    `gorm:"primaryKey"`
	AccountID string    `gorm:"not null"`
	Account   Account   `gorm:"foreignKey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	RefCode   string    `gorm:"not null"`
	Amount    float64   `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt
}

func (Payment) TableName() string {
	return "payments"
}

func (p *Payment) BeforeCreate(tx *gorm.DB) (err error) {

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

	p.ID = id
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return
}

func (p *Payment) BeforeUpdate(tx *gorm.DB) (err error) {
	p.UpdatedAt = time.Now()
	return
}
