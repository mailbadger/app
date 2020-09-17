package entities

import "time"

// Click entity holds information regarding link clicks.
type Click struct {
	ID         int64     `json:"id" gorm:"column:id; primary_key:yes"`
	UserID     int64     `json:"-"`
	CampaignID int64     `json:"campaign_id"`
	Recipient  string    `json:"recipient"`
	Link       string    `json:"link"`
	UserAgent  string    `json:"user_agent"`
	IPAddress  string    `json:"ip_address"`
	CreatedAt  time.Time `json:"created_at"`
}

// ClicksStats entity holds information about total and unique recipients that clicked a certain link
type ClicksStats struct {
	Link         string `json:"link,omitempty"`
	UniqueClicks int64  `json:"unique_clicks"`
	TotalClicks  int64  `json:"total_clicks"`
}
