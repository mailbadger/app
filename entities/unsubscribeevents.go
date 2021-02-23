package entities

import (
	"time"

	"github.com/segmentio/ksuid"
)

type UnsubscribeEvents struct {
	ID        ksuid.KSUID `json:"id" gorm:"column:id; primary_key:yes"`
	Email     string      `json:"email"`
	CreatedAt time.Time   `json:"created_at"`
}
