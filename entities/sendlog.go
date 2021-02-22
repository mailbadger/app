package entities

import (
	"time"

	"github.com/segmentio/ksuid"
)

const (
	// SendLogStatusFailed status used when sender consumer fails to send email
	SendLogStatusFailed = "failed"
	// SendLogStatusSuccessful status used when sender consumer succeeded sending the mail
	SendLogStatusSuccessful = "successful"

	// SendLogDescriptionOnSuccessful description used when sender succeeded sending the mail
	SendLogDescriptionOnSuccessful = "Email sent successfully"
	// SendLogDescriptionOnSesClientError description used when sender failed to create ses client
	SendLogDescriptionOnSesClientError = "Unable to create ses client"
	// SendLogDescriptionOnSendEmailError description used when ses client fails to send the email
	SendLogDescriptionOnSendEmailError = "Unable to send email"
)

type SendLog struct {
	ID           ksuid.KSUID `json:"id" gorm:"column:id; primary_key:yes"`
	MessageID    *string     `json:"message_id" gorm:"message_id; index"`
	UserID       int64       `json:"-" gorm:"column:user_id; index"`
	SubscriberID int64       `json:"-" gorm:"column:subscriber_id"`
	CampaignID   int64       `json:"campaign_id" gorm:"column:campaign_id; index"`
	Status       string      `json:"status"`
	Description  string      `json:"description"`
	CreatedAt    time.Time   `json:"created_at"`
}
