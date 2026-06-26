package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	RoleID uuid.UUID `gorm:"type:uuid;primary_key"`
	Name string `gorm:"type:varchar(50);unique;not null"` // "superadmin", "admin", "user"
}

func (r *Role) BeforeCreate(tx *gorm.DB) (err error) {
	if r.RoleID == uuid.Nil {
		r.RoleID = uuid.New()
	}
	return
}