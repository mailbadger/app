package entities

import (
	"time"

	"github.com/segmentio/ksuid"
)

type UnsubscribeEvent struct {
	ID        ksuid.KSUID `json:"id" gorm:"column:id; primary_key:yes"`
	UserID    int64       `json:"user_id"`
	Email     string      `json:"email"`
	CreatedAt time.Time   `json:"created_at"`
}
