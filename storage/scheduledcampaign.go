package storage

import (
	"github.com/mailbadger/app/entities"
)

// CreateScheduledCampaign creates a report.
func (db *store) CreateScheduledCampaign(c *entities.ScheduledCampaign) error {
	return db.Create(c).Error
}

// DeleteScheduledCampaign deletes scheduled campaign.
func (db *store) DeleteScheduledCampaign(campaignID int64) error {
	return db.Where("campaign_id = ?", campaignID).Delete(entities.ScheduledCampaign{}).Error
}
