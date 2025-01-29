package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title    string `json:"title" gorm:"not null"`
	Content  string `json:"content" gorm:"type:text;not null"`
	AuthorID uint   `json:"author_id"`
	Author   User   `json:"author" gorm:"foreignKey:AuthorID"`
	Likes    uint   `json:"likes" gorm:"default:0"`
}
