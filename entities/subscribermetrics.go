package entities

import (
	"time"
)

// SubscriberMetrics represents daily events per user
type SubscriberMetrics struct {
	ID           int64     `json:"id" gorm:"column:id"`
	UserID       int64     `json:"user_id" gorm:"primaryKey"`
	Created      int64     `json:"created"`
	Unsubscribed int64     `json:"unsubscribed"`
	Date         time.Time `json:"date" gorm:"primaryKey"`
}
