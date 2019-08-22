package entities

import (
	"encoding/json"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/sirupsen/logrus"
)

//Subscriber represents the subscriber entity
type Subscriber struct {
	ID          int64             `json:"id" gorm:"column:id; primary_key:yes"`
	UserID      int64             `json:"-" gorm:"column:user_id; index"`
	Name        string            `json:"name" gorm:"not null"`
	Email       string            `json:"email" gorm:"not null"`
	MetaJSON    JSON              `json:"-" gorm:"column:metadata; type:json"`
	Segments    []Segment         `json:"-" gorm:"many2many:subscribers_segments;"`
	Blacklisted bool              `json:"blacklisted"`
	Active      bool              `json:"active"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Errors      map[string]string `json:"-" sql:"-"`
	Metadata    map[string]string `json:"metadata" sql:"-"`
}

func (s *Subscriber) Normalize() {
	m := make(map[string]string)

	if !s.MetaJSON.IsNull() {
		err := json.Unmarshal(s.MetaJSON, &m)
		if err != nil {
			logrus.WithError(err).Error("unable to unmarshal json metadata")
		}
	}

	m["name"] = s.Name

	s.Metadata = m
}

// Validate subscriber properties,
func (s *Subscriber) Validate() bool {
	s.Errors = make(map[string]string)

	if valid.Trim(s.Name, "") == "" {
		s.Errors["name"] = "The subscriber name cannot be empty."
	}

	if !valid.IsEmail(s.Email) {
		s.Errors["email"] = "The specified email is not valid."
	}

	return len(s.Errors) == 0
}
