package entities

import "time"

//Campaign represents the campaign entity
type Campaign struct {
	Id          int64     `json:"id"`
	UserId      int64     `json:"-"`
	Name        string    `json:"name" sql:"not null"`
	Subject     string    `json:"subject"`
	TemplateId  int64     `json:"-"`
	Template    Template  `json:"template"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	ScheduledAt time.Time `json:"scheduled_at"`
	CompletedAt time.Time `json:"completed_at"`
}
