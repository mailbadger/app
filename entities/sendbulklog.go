package entities

import "time"

type SendBulkLog struct {
	ID         int64     `json:"id" gorm:"column:id; primary_key:yes"`
	UUID       string    `json:"uuid" gorm:"column:uuid; index"`
	UserID     int64     `json:"-" gorm:"column:user_id; index"`
	CampaignID int64     `json:"campaign_id" gorm:"column:campaign_id; index"`
	MessageID  string    `json:"message_id"`
	Status     string    `json:"status"`
	Error      *string   `json:"error"`
	CreatedAt  time.Time `json:"created_at"`
}
