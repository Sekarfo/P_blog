package subscription

import (
	"time"

	"gorm.io/gorm"
)

type Subscription struct {
	gorm.Model
	UserID            uint       `json:"user_id"`
	Status            string     `json:"status" gorm:"default:'pending'"` // pending, approved, rejected
	RequestedAt       time.Time  `json:"requested_at"`
	ApprovedAt        *time.Time `json:"approved_at,omitempty"`
	SubscriptionStart *time.Time `json:"subscription_start,omitempty"`
	SubscriptionEnd   *time.Time `json:"subscription_end,omitempty"`
}
