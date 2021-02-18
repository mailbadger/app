package entities

import (
	"time"
)

type UnsubscribeEvents struct {
	ID        int64     `json:"id" gorm:"column:id; primary_key:yes"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
