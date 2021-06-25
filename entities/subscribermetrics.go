package entities

import (
	"time"
)

// SubscribersMetrics represents daily events per user
type SubscribersMetrics struct {
	ID           int64     `json:"id" gorm:"column:id"`
	UserID       int64     `json:"user_id" gorm:"primaryKey"`
	Created      int64     `json:"created"`
	Deleted      int64     `json:"deleted"`
	Unsubscribed int64     `json:"unsubscribed"`
	Date         time.Time `json:"date" gorm:"primaryKey"`
}
