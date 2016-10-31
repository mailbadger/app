package entities

import (
	"errors"
	"time"

	valid "github.com/asaskevich/govalidator"
)

const (
	STATUS_DRAFT     = "draft"
	STATUS_COMPLETED = "completed"
	STATUS_SENDING   = "sending"
	STATUS_SCHEDULED = "scheduled"
)

//Campaign represents the campaign entity
type Campaign struct {
	Id          int64             `json:"id"`
	UserId      int64             `json:"-" gorm:"column:user_id; index"`
	Name        string            `json:"name" gorm:"not null"`
	Subject     string            `json:"subject"`
	TemplateId  int64             `json:"-" gorm:"column:template_id; index"`
	Template    Template          `json:"template"`
	Status      string            `json:"status"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	ScheduledAt time.Time         `json:"scheduled_at"`
	CompletedAt time.Time         `json:"completed_at"`
	Errors      map[string]string `json:"-" sql:"-"`
}

var ErrCampaignNameEmpty = errors.New("The name cannot be empty.")
var ErrSubjectEmpty = errors.New("The subject cannot be empty.")
var ErrTemplateNotSpecified = errors.New("The template id must be specified.")

// Validate campaign properties,
func (c *Campaign) Validate() bool {
	c.Errors = make(map[string]string)

	if valid.Trim(c.Name, "") == "" {
		c.Errors["name"] = ErrTemplateNameEmpty.Error()
	}
	if valid.Trim(c.Subject, "") == "" {
		c.Errors["subject"] = ErrSubjectEmpty.Error()
	}

	return len(c.Errors) == 0
}
