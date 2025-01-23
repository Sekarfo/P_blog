package models

import "time"

type Blog struct {
	ID        uint      `gorm:"primaryKey" json:"blog_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorID  uint      `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
}
