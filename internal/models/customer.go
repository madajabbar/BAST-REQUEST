package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	CustomerID   uuid.UUID `gorm:"type:uuid;primary_key"`
	CustomerCode string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	CustomerName string    `gorm:"type:varchar(255);not null"`
	Status       string    `gorm:"type:varchar(50);default:'active'"` // active, inactive
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// TableName overrides the table name used by Customer to `master_customer`
func (Customer) TableName() string {
	return "master_customer"
}

func (c *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	if c.CustomerID == uuid.Nil {
		c.CustomerID = uuid.New()
	}
	return
}
