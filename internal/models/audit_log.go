package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AuditLog struct {
	AuditLogID  uuid.UUID      `gorm:"type:uuid;primary_key"`
	TargetTable string         `gorm:"type:varchar(100);not null;column:table_name"`
	RecordID    string         `gorm:"type:varchar(100);not null"`
	Action      string         `gorm:"type:varchar(20);not null"` // POST, PUT, DELETE, PATCH
	OldData     datatypes.JSON `gorm:"type:jsonb"`
	NewData     datatypes.JSON `gorm:"type:jsonb"`
	PerformedBy string         `gorm:"type:varchar(100);not null"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
}

func (AuditLog) TableName() string {
	return "audit_log"
}

func (a *AuditLog) BeforeCreate(tx *gorm.DB) (err error) {
	if a.AuditLogID == uuid.Nil {
		a.AuditLogID = uuid.New()
	}
	return
}
