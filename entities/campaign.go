package entities

import (
	"time"

	"github.com/aws/aws-sdk-go/service/ses"
)

const (
	// StatusDraft indicates a draft campaign ready to be sent.
	StatusDraft = "draft"
	// StatusSending indicates that the campaign is in the sending process.
	StatusSending = "sending"
	// StatusSent indicates that a campaign has been sent.
	StatusSent = "sent"
	// StatusScheduled indicates a scheduled campaign status.
	StatusScheduled = "scheduled"
	// CampaignsTopic is the topic used by the campaigner consumer.
	CampaignsTopic = "campaigns"
	// SendBulkTopic is the topic used by the bulksender consumer.
	SendBulkTopic = "send_bulk"
)

//Campaign represents the campaign entity
type Campaign struct {
	Model
	UserID      int64             `json:"-" gorm:"column:user_id; index"`
	Name        string            `json:"name" gorm:"not null"`
	TemplateID  int64             `json:"-"`
	Template    *Template         `json:"template" gorm:"foreignKey:template_id"`
	Status      string            `json:"status"`
	ScheduledAt NullTime          `json:"scheduled_at" gorm:"column:scheduled_at"`
	CompletedAt NullTime          `json:"completed_at" gorm:"column:completed_at"`
	DeletedAt   NullTime          `json:"deleted_at" gorm:"column:deleted_at"`
	Errors      map[string]string `json:"-" sql:"-"`
}

// BulkSendMessage represents the entity used to transport the bulk send message
// used by the bulksender consumer.
type BulkSendMessage struct {
	UUID       string                           `json:"msg_uuid"`
	UserID     int64                            `json:"user_id"`
	CampaignID int64                            `json:"campaign_id"`
	SesKeys    *SesKeys                         `json:"ses_keys"`
	Input      *ses.SendBulkTemplatedEmailInput `json:"input"`
}

// SendCampaignParams represent the request params used
// by the send campaign endpoint.
type SendCampaignParams struct {
	SegmentIDs             []int64           `json:"segment_ids"`
	TemplateData           map[string]string `json:"template_data"`
	Source                 string            `json:"source"`
	UserID                 int64             `json:"user_id"`
	UserUUID               string            `json:"user_uuid"`
	ConfigurationSetExists bool              `json:"configuration_set_exists"`
	Campaign               `json:"campaign"`
	SesKeys                `json:"ses_keys"`
}

// CampaignClicksStats represents clicks stats by campaign, total number of links and stats for each link
type CampaignClicksStats struct {
	Total       int64         `json:"total"`
	ClicksStats []ClicksStats `json:"collection"`
}

func (c Campaign) GetID() int64 {
	return c.Model.ID
}

func (c Campaign) GetCreatedAt() time.Time {
	return c.Model.CreatedAt
}

func (c Campaign) GetUpdatedAt() time.Time {
	return c.Model.UpdatedAt
}

type OpensStats struct {
	Unique int64 `json:"unique"`
	Total  int64 `json:"total"`
}

type CampaignStats struct {
	TotalSent  int64        `json:"total_sent"`
	Delivered  int64        `json:"delivered"`
	Opens      *OpensStats  `json:"opens"`
	Clicks     *ClicksStats `json:"clicks"`
	Bounces    int64        `json:"bounces"`
	Complaints int64        `json:"complaints"`
}
