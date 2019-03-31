package entities

import "time"

type Complaint struct {
	ID         int64     `json:"id" gorm:"column:id; primary_key:yes"`
	UserID     int64     `json:"-"`
	CampaignID int64     `json:"campaign_id"`
	Recipient  string    `json:"recipient"`
	UserAgent  string    `json:"user_agent"`
	Type       string    `json:"type"`
	FeedbackID string    `json:"feedback_id"`
	CreatedAt  time.Time `json:"created_at"`
}
