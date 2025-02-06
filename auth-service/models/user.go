package models

import "time"

type Role struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID              uint      `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Password        string    `json:"-"`
	RoleID          uint      `json:"role_id"`
	Role            Role      `json:"role"`
	IsEmailVerified bool      `json:"is_email_verified"`
	ResetToken      string    `json:"-"` // Add ResetToken field
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
