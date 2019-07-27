package entities

import "time"

type AccessKey struct {
	ID        int64     `json:"id" gorm:"column:id; primary_key:yes"`
	UserID    int64     `json:"-"`
	User      User      `json:"-"`
	AccessKey string    `json:"access_key"`
	SecretKey string    `json:"secret_key"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
