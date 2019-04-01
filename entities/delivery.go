package entities

import "time"

type Delivery struct {
	ID                   int64     `json:"id" gorm:"column:id; primary_key:yes"`
	UserID               int64     `json:"-"`
	CampaignID           int64     `json:"campaign_id"`
	Recipient            string    `json:"recipient"`
	ProcessingTimeMillis int64     `json:"processing_time_millis"`
	SMTPResponse         string    `json:"smtp_response"`
	ReportingMTA         string    `json:"reporting_mta"`
	RemoteMtaIP          string    `json:"remote_mta_ip"`
	CreatedAt            time.Time `json:"created_at"`
}
