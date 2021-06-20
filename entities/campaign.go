package entities

import (
	"time"

	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/cbroglie/mustache"
	"github.com/segmentio/ksuid"
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
	// CampaignerTopic is the topic used by the campaigner consumer.
	CampaignerTopic = "campaigner"
	// SendBulkTopic is the topic used by the bulksender consumer.
	SendBulkTopic = "send_bulk"
	// SenderTopic is the topic used by the sender consumer.
	SenderTopic = "sender"
)

// Campaign represents the campaign entity
type Campaign struct {
	Model
	UserID       int64             `json:"-" gorm:"column:user_id; index"`
	EventID      *ksuid.KSUID      `json:"-"`
	Name         string            `json:"name" gorm:"not null"`
	TemplateID   int64             `json:"-"`
	BaseTemplate *BaseTemplate     `json:"template" gorm:"foreignKey:template_id"`
	Schedule     *CampaignSchedule `json:"schedule" gorm:"foreignKey:campaign_id"`
	Status       string            `json:"status"`
	CompletedAt  NullTime          `json:"completed_at" gorm:"column:completed_at"`
	DeletedAt    NullTime          `json:"-" gorm:"column:deleted_at"`
	StartedAt    NullTime          `json:"started_at" gorm:"column:started_at"`
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

// CampaignerTopicParams represent the request params used
// by the send campaign endpoint.
type CampaignerTopicParams struct {
	EventID                ksuid.KSUID       `json:"event_id"`
	CampaignID             int64             `json:"campaign_id"`
	SegmentIDs             []int64           `json:"segment_ids"`
	TemplateData           map[string]string `json:"template_data"`
	Source                 string            `json:"source"`
	UserID                 int64             `json:"user_id"`
	UserUUID               string            `json:"user_uuid"`
	ConfigurationSetExists bool              `json:"configuration_set_exists"`
	SesKeys                `json:"ses_keys"`
}

// SenderTopicParams represent the request params used
// by the sender campaign consumer.
type SenderTopicParams struct {
	EventID                ksuid.KSUID `json:"event_id"`
	UserID                 int64       `json:"user_id"`
	UserUUID               string      `json:"user_uuid"`
	CampaignID             int64       `json:"campaign_id"`
	SubscriberID           int64       `json:"subscriber_id"`
	SubscriberEmail        string      `json:"subscriber_email"`
	Source                 string      `json:"source"`
	ConfigurationSetExists bool        `json:"configuration_set_exists"`
	HTMLPart               []byte      `json:"html_part"`
	SubjectPart            []byte      `json:"subject_part"`
	TextPart               []byte      `json:"text_part"`
	SesKeys                SesKeys     `json:"ses_keys"`
}

type CampaignTemplateData struct {
	Template    *Template
	HTMLPart    *mustache.Template
	SubjectPart *mustache.Template
	TextPart    *mustache.Template
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

// SetCampaignEventID if the campaign is scheduled then sets the id to the scheduled campaign's id else generates new id
func (c *Campaign) SetEventID() {
	if c.Schedule != nil {
		c.EventID = &c.Schedule.ID
		return
	}

	uid := ksuid.New()
	c.EventID = &uid
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
