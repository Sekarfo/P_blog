package models

import "gorm.io/gorm"

type Like struct {
	gorm.Model
	PostID uint `json:"post_id"`
	Post   Post `json:"post" gorm:"foreignKey:PostID"`
	UserID uint `json:"user_id"`
	User   User `json:"user" gorm:"foreignKey:UserID"`
}
