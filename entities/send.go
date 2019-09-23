package entities

import "time"

type Send struct {
	ID               int64     `json:"id" gorm:"column:id; primary_key:yes"`
	UserID           int64     `json:"-" gorm:"column:user_id; index"`
	CampaignID       int64     `json:"campaign_id"`
	MessageID        string    `json:"message_id"`
	Source           string    `json:"source"`
	SendingAccountID string    `json:"sending_account_id"`
	Destination      string    `json:"destination"`
	CreatedAt        time.Time `json:"created_at"`
}
