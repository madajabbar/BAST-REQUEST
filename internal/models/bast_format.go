package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BastFormat struct {
	FormatID      uuid.UUID `gorm:"type:uuid;primary_key"`
	FormatName    string    `gorm:"type:varchar(100);not null"`
	FormatType    string    `gorm:"type:varchar(50);not null"` // PO, Internal
	FormatPattern string    `gorm:"type:varchar(255);not null"` // e.g. BAST/{YYYY}/{MM}/{SEQ}
	IsActive      bool      `gorm:"default:true"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func (BastFormat) TableName() string {
	return "master_bast_format"
}

func (b *BastFormat) BeforeCreate(tx *gorm.DB) (err error) {
	if b.FormatID == uuid.Nil {
		b.FormatID = uuid.New()
	}
	return
}
