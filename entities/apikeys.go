package entities

import "time"

// APIKey represents the user key used to authenticate requests
// with the API.
type APIKey struct {
	ID        int64     `json:"id" gorm:"column:id; primary_key:yes"`
	UserID    int64     `json:"-"`
	User      User      `json:"-"`
	SecretKey string    `json:"secret_key"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
