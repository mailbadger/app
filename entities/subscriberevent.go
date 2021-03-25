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
