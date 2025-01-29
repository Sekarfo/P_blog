package models

import (
	"gorm.io/gorm"
)

// Role model
type Role struct {
	gorm.Model
	Name string `gorm:"unique;not null"`
}

// User model
type User struct {
	gorm.Model
	Name            string `gorm:"not null"`
	Email           string `gorm:"unique;not null"`
	Password        string `gorm:"not null"`
	RoleID          uint
	Role            Role
	IsEmailVerified bool `json:"is_email_verified" gorm:"default:false"`
}
