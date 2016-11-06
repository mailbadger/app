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
	Metadata    []ListMetadata    `json:"metadata" gorm:"ForeignKey:ListId"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Errors      map[string]string `json:"-" sql:"-"`
}

//ListMetadata represents the list metadata in a form of a key and value
type ListMetadata struct {
	Id        int64  `gorm:"column:id; primary_key:yes"`
	ListId    int64  `gorm:"column:list_id; index"`
	Key       string `gorm:"not null"`
	Value     string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Validate campaign properties,
func (c *List) Validate() bool {
	c.Errors = make(map[string]string)

	if valid.Trim(c.Name, "") == "" {
		c.Errors["name"] = ErrListNameEmpty.Error()
	}

	return len(c.Errors) == 0
}
