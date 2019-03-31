package entities

import "time"

type Click struct {
	ID         int64     `json:"id" gorm:"column:id; primary_key:yes"`
	UserID     int64     `json:"-"`
	CampaignID int64     `json:"campaign_id"`
	Link       string    `json:"link"`
	UserAgent  string    `json:"user_agent"`
	IPAddress  string    `json:"ip_address"`
	CreatedAt  time.Time `json:"created_at"`
}
