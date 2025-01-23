package models

import "time"

type Commentary struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	BlogID    uint      `json:"blog_id"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}
