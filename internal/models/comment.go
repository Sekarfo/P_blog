package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Content  string `json:"content" gorm:"not null"`
	PostID   uint   `json:"post_id"`
	Post     Post   `json:"post" gorm:"foreignKey:PostID"`
	AuthorID uint   `json:"author_id"`
	Author   User   `json:"author" gorm:"foreignKey:AuthorID"`
}
