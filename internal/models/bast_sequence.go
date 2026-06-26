package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BastSequence struct {
	SequenceID uuid.UUID  `gorm:"type:uuid;primary_key"`
	FormatID   uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:idx_format_year_month"`
	Format     BastFormat `gorm:"foreignKey:FormatID;references:FormatID"`
	Year       int        `gorm:"not null;uniqueIndex:idx_format_year_month"`
	Month      int        `gorm:"not null;uniqueIndex:idx_format_year_month"`
	LastNumber int        `gorm:"not null;default:0"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (BastSequence) TableName() string {
	return "bast_sequence"
}

func (b *BastSequence) BeforeCreate(tx *gorm.DB) (err error) {
	if b.SequenceID == uuid.Nil {
		b.SequenceID = uuid.New()
	}
	return
}
