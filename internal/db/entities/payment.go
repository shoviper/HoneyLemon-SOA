package entities

import (
	"time"

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

func (c *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return
}

func (c *Payment) BeforeUpdate(tx *gorm.DB) (err error) {
	c.UpdatedAt = time.Now()
	return
}
