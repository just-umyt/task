package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User struct to describe User object.
type User struct {
	ID               uuid.UUID `gorm:"type:uuid" validate:"required,uuid"`
	Email            string    `gorm:"unique;type:varchar(255)" validate:"required,email,lte=255"`
	PasswordHash     string    `gorm:"type:varchar(255)" validate:"required,lte=255"`
	RefreshTokenHash string    `gorm:"type:varchar(255)"`
	gorm.Model
}
