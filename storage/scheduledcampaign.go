package storage

import (
	"github.com/segmentio/ksuid"

	"github.com/mailbadger/app/entities"
)

// CreateScheduledCampaign creates a report.
func (db *store) CreateScheduledCampaign(c *entities.ScheduledCampaign) error {
	return db.Create(c).Error
}

// DeleteScheduledCampaign creates a report.
func (db *store) DeleteScheduledCampaign(id ksuid.KSUID, campaignID int64) error {
	return db.Where("campaign_id = ?", campaignID).Delete(entities.ScheduledCampaign{ID: id}).Error
}
