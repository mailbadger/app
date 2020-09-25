package entities

import "time"

type Model struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ModelInterface interface {
	GetID() int64
}
