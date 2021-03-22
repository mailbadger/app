package entities

import (
	"time"

	"github.com/segmentio/ksuid"
)

type CampaignFailedLog struct {
	ID          ksuid.KSUID `json:"id" gorm:"column:id; primary_key:yes"`
	UserID      int64       `json:"user_id"`
	CampaignID  int64       `json:"campaign_id"`
	Description string      `json:"description"`
	CreatedAt   time.Time   `json:"created_at"`
}
