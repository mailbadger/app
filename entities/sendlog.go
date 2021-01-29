package entities

import "time"

type SendLog struct {
	ID           int64     `json:"id" gorm:"column:id; primary_key:yes"`
	UUID         string    `json:"uuid" gorm:"column:uuid; index"`
	UserID       int64     `json:"-" gorm:"column:user_id; index"`
	SubscriberID int64     `json:"-" gorm:"column:subscriber_id"`
	CampaignID   int64     `json:"campaign_id" gorm:"column:campaign_id; index"`
	Status       string    `json:"status"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
}

const (
	FailedSendLogStatus     = "failed"
	SuccessfulSendLogStatus = "successful"
)
