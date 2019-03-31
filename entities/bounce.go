package entities

import "time"

type Bounce struct {
	ID             int64     `json:"id" gorm:"column:id; primary_key:yes"`
	UserID         int64     `json:"-"`
	CampaignID     int64     `json:"campaign_id"`
	Recipient      string    `json:"recipient"`
	Type           string    `json:"type"`
	SubType        string    `json:"sub_type"`
	Action         string    `json:"action"`
	Status         string    `json:"status"`
	DiagnosticCode string    `json:"diagnostic_code"`
	FeedbackID     string    `json:"feedback_id"`
	CreatedAt      time.Time `json:"created_at"`
}
