package entities

import (
	"time"
)

// SesKeys entity holds information about the client's
// SES access and secret key.
type SesKeys struct {
	ID        int64     `json:"id" gorm:"column:id; primary_key:yes"`
	UserID    int64     `json:"-" gorm:"column:user_id; index"`
	AccessKey string    `json:"access_key" gorm:"not null" valid:"alphanum,required"`
	SecretKey string    `json:"secret_key,omitempty" gorm:"not null" valid:"required"`
	Region    string    `json:"region" gorm:"not null" valid:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
