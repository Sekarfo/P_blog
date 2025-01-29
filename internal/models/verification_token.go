package models

import (
	"time"

	"gorm.io/gorm"
)

type VerificationToken struct {
	gorm.Model
	UserID    uint      `json:"user_id"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	Token     string    `json:"token" gorm:"unique;not null"`
	ExpiresAt time.Time `json:"expires_at"`
}
