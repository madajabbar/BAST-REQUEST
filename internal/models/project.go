package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	ProjectID   uuid.UUID `gorm:"type:uuid;primary_key"`
	CustomerID  uuid.UUID `gorm:"type:uuid;not null"`
	Customer    Customer  `gorm:"foreignKey:CustomerID;references:CustomerID"`
	ProjectCode string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	ProjectName string    `gorm:"type:varchar(255);not null"`
	Status      string    `gorm:"type:varchar(50);default:'active'"` // active, inactive
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (Project) TableName() string {
	return "master_project"
}

func (p *Project) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ProjectID == uuid.Nil {
		p.ProjectID = uuid.New()
	}
	return
}
