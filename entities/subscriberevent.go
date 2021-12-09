package entities

import (
	"time"
	
	"github.com/segmentio/ksuid"
)

const (
	SubscriberEventTypeCreated      EventType = "created"
	SubscriberEventTypeDeleted      EventType = "deleted"
	SubscriberEventTypeUnsubscribed EventType = "unsubscribed"
)

// SubscriberEvent represents an event saved on subscriber's change
type SubscriberEvent struct {
	ID              ksuid.KSUID `json:"id" gorm:"column:id; primary_key:yes"`
	UserID          int64       `json:"user_id"`
	SubscriberEmail string      `json:"subscriber_email"`
	EventType       EventType   `json:"event_type"`
	CreatedAt       time.Time   `json:"created_at"`
}

// GroupedSubscriberEvents represents grouped subscriber events records by user date and event type
type GroupedSubscriberEvents struct {
	UserID    int64     `json:"user_id"`
	Date      time.Time `json:"date"`
	EventType string `json:"event_type"`
	Total     int64     `json:"total"`
}
