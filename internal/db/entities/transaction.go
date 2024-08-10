package entities

import (
	"time"

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

func (c *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return
}

func (c *Transaction) BeforeUpdate(tx *gorm.DB) (err error) {
	c.UpdatedAt = time.Now()
	return
}
