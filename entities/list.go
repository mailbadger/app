package entities

import (
	"errors"
	"time"

	valid "github.com/asaskevich/govalidator"
)

var ErrListNameEmpty = errors.New("The list name cannot be empty.")

//List represents the list entity
type List struct {
	Id          int64             `json:"id"`
	Name        string            `json:"name" gorm:"not null"`
	UserId      int64             `json:"-" gorm:"column:user_id; index"`
	Subscribers []Subscriber      `gorm:"many2many:subscribers_lists;"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Errors      map[string]string `json:"-" sql:"-"`
}

// Validate list properties,
func (l *List) Validate() bool {
	l.Errors = make(map[string]string)

	if valid.Trim(l.Name, "") == "" {
		l.Errors["name"] = ErrListNameEmpty.Error()
	}

	return len(l.Errors) == 0
}
