package entities

import (
	"time"
    cryptoRand "crypto/rand"
    "math/big"

    "github.com/oklog/ulid/v2"

	"gorm.io/gorm"
)

type Account struct {
	ID        string    `gorm:"primaryKey"`
	ClientID  string    `gorm:"not null"`
	Client    Client    `gorm:"foreignKey:ClientID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Type      string    `gorm:"not null"`
	Balance   float64   `gorm:"not null"`
	Pin       string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt
}

func (Account) TableName() string {
	return "accounts"
}

func (c *Account) BeforeCreate(tx *gorm.DB) (err error) {
	// Create a new ULID
	t := time.Now().UTC()
	entropy := ulid.Monotonic(cryptoRand.Reader, 0)
	ulid := ulid.MustNew(ulid.Timestamp(t), entropy)

	// Convert the ULID to a big.Int
	bigInt := new(big.Int).SetBytes(ulid[:])

	// Get the last 10 digits
	modulus := new(big.Int).SetInt64(10000000000)
	result := new(big.Int).Mod(bigInt, modulus)

	// Convert the result to a string
	c.ID = result.String()
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return
}

func (c *Account) BeforeUpdate(tx *gorm.DB) (err error) {
	c.UpdatedAt = time.Now()
	return
}
