package entities

import (
	"time"

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
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return
}

func (c *Account) BeforeUpdate(tx *gorm.DB) (err error) {
	c.UpdatedAt = time.Now()
	return
}