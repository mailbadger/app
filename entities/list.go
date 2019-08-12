package entities

import (
	"time"

	valid "github.com/asaskevich/govalidator"
)

// List represents the list entity
type List struct {
	ID          int64             `json:"id"`
	Name        string            `json:"name" gorm:"not null"`
	UserID      int64             `json:"-" gorm:"column:user_id; index"`
	Subscribers []Subscriber      `json:"-" gorm:"many2many:subscribers_lists;"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Errors      map[string]string `json:"-" sql:"-"`
}

// Validate validates the list properties and populates the Errors map
// in case of any errors.
func (l *List) Validate() bool {
	l.Errors = make(map[string]string)

	if valid.Trim(l.Name, "") == "" {
		l.Errors["name"] = "The list name cannot be empty."
	}

	return len(l.Errors) == 0
}
