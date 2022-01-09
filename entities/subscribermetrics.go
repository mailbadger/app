package entities

import (
	"time"
)

// SubscriberMetrics represents daily events per user
type SubscriberMetrics struct {
	UserID       int64     `json:"user_id"`
	Created      int64     `json:"created"`
	Unsubscribed int64     `json:"unsubscribed"`
	Datetime     time.Time `json:"datetime"`
}
