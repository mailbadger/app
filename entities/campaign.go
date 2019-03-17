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
	Id           int64             `json:"id" gorm:"column:id; primary_key:yes"`
	UserId       int64             `json:"-" gorm:"column:user_id; index"`
	Name         string            `json:"name" gorm:"not null" valid:"alphanum,required"`
	TemplateName string            `json:"template_name" valid:"alphanum,required"`
	Status       string            `json:"status"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	ScheduledAt  Time              `json:"scheduled_at"`
	CompletedAt  Time              `json:"completed_at"`
	Errors       map[string]string `json:"-" sql:"-"`
}

var ErrCampaignNameEmpty = errors.New("The name cannot be empty.")

// Validate campaign properties,
func (c *Campaign) Validate() bool {
	c.Errors = make(map[string]string)

	if valid.Trim(c.Name, "") == "" {
		c.Errors["name"] = ErrCampaignNameEmpty.Error()
	}

	res, err := valid.ValidateStruct(c)
	if err != nil || !res {
		c.Errors["reason"] = err.Error()
	}

	return len(c.Errors) == 0
}
