package entities

import "time"

//Subscriber represents the subscriber entity
type Subscriber struct {
	Id        int64             `json:"id"`
	Name      string            `json:"name" gorm:"not null"`
	Email     string            `json:"email" gorm:"not null"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Errors    map[string]string `json:"-" sql:"-"`
}
