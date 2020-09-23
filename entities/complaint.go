package entities

import "time"

// Complaint represents an entity regarding user complaint information.
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

func (c Complaint) GetCreatedAt() time.Time {
	return c.CreatedAt
}

func (c Complaint) GetID() int64 {
	return c.ID
}

func (c Complaint) GetUpdatedAt() time.Time {
	return time.Time{}
}
