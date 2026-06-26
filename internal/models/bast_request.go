package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BastRequest struct {
	BastRequestID uuid.UUID  `gorm:"type:uuid;primary_key"`
	CustomerID    uuid.UUID  `gorm:"type:uuid;not null"`
	Customer      Customer   `gorm:"foreignKey:CustomerID;references:CustomerID"`
	ProjectID     uuid.UUID  `gorm:"type:uuid;not null"`
	Project       Project    `gorm:"foreignKey:ProjectID;references:ProjectID"`
	FormatID      uuid.UUID  `gorm:"type:uuid;not null"`
	Format        BastFormat `gorm:"foreignKey:FormatID;references:FormatID"`
	Perihal       string     `gorm:"type:varchar(255);not null"`
	TipeNomor     string     `gorm:"type:varchar(50);not null"` // PO, Internal
	PoNumber      string     `gorm:"type:varchar(100)"`
	BastNumber    string     `gorm:"type:varchar(100);not null;uniqueIndex"`
	Status        string     `gorm:"type:varchar(50);default:'Active'"` // Active, Used, Void
	RequestedBy   string     `gorm:"type:varchar(100);not null"`
	RequestedAt   time.Time  `gorm:"autoCreateTime"`
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func (BastRequest) TableName() string {
	return "bast_request"
}

func (b *BastRequest) BeforeCreate(tx *gorm.DB) (err error) {
	if b.BastRequestID == uuid.Nil {
		b.BastRequestID = uuid.New()
	}
	return
}
