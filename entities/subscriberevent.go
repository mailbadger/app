package entities

import (
	"time"
)

const (
	SubscriberEventTypeCreated      EventType = "created"
	SubscriberEventTypeUnsubscribed EventType = "unsubscribed"
)

// SubscriberEvent represents an event saved on subscriber's change
type SubscriberEvent struct {
	ID           int64     `json:"id" gorm:"column:id; primary_key:yes"`
	UserID       int64     `json:"user_id"`
	SubscriberID int64     `json:"subscriber_id"`
	EventType    EventType `json:"event_type"`
	CreatedAt    time.Time `json:"created_at"`
}

// GroupedSubscriberEvents represents grouped subscriber events records by user date and event type
type GroupedSubscriberEvents struct {
	UserID    int64     `json:"user_id"`
	Date      time.Time `json:"date"`
	EventType string    `json:"event_type"`
	Total     int64     `json:"total"`
}
