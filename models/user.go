package models

import "time"

type User struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	Name              string    `json:"name"`
	Email             string    `json:"email"`
	Password          string    `json:"password"`
	EmailVerified     bool      `json:"email_verified"`
	EmailVerifyToken  string    `json:"-"`
	EmailVerifyExpiry time.Time `json:"-"`
}
