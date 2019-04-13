package entities

import (
	"errors"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/aws/aws-sdk-go/service/ses"
)

const (
	StatusDraft     = "draft"
	StatusSending   = "sending"
	StatusSent      = "sent"
	StatusScheduled = "scheduled"

	CampaignsTopic = "campaigns"
	SendBulkTopic  = "send_bulk"
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
	ScheduledAt  NullTime          `json:"scheduled_at" gorm:"column:scheduled_at"`
	CompletedAt  NullTime          `json:"completed_at" gorm:"column:completed_at"`
	Errors       map[string]string `json:"-" sql:"-"`
}

type BulkSendMessage struct {
	UUID       string                           `json:"msg_uuid"`
	UserID     int64                            `json:"user_id"`
	CampaignID int64                            `json:"campaign_id"`
	SesKeys    *SesKeys                         `json:"ses_keys"`
	Input      *ses.SendBulkTemplatedEmailInput `json:"input"`
}

type SendCampaignParams struct {
	ListIDs      []int64           `json:"list_ids"`
	TemplateData map[string]string `json:"template_data"`
	Source       string            `json:"source"`
	UserID       int64             `json:"user_id"`
	Campaign     `json:"campaign"`
	SesKeys      `json:"ses_keys"`
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
		c.Errors["message"] = err.Error()
	}

	return len(c.Errors) == 0
}
