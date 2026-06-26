package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	UserID uuid.UUID `gorm:"type:uuid;primary_key"`
	RoleID uuid.UUID `gorm:"type:uuid;not null"`
	Role   Role      `gorm:"foreignKey:RoleID"`
	Username string  `gorm:"type:varchar(100);not null"`
	Email string `gorm:"type:varchar(100);uniqueIndex;not null"`
	PasswordHash string `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.UserID == uuid.Nil {
		u.UserID = uuid.New()
	}
	return
}
