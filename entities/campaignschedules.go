package entities

import (
	"time"

	"github.com/segmentio/ksuid"
)

type CampaignSchedule struct {
	ID          ksuid.KSUID `json:"id" gorm:"column:id; primary_key:yes"`
	UserID      int64       `json:"user_id"`
	CampaignID  int64       `json:"campaign_id"`
	ScheduledAt time.Time   `json:"scheduled_at"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}
