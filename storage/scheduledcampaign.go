package storage

import (
	"errors"

	"github.com/jinzhu/gorm"

	"github.com/mailbadger/app/entities"
)

// CreateScheduledCampaign creates a scheduled campaign.
func (db *store) CreateScheduledCampaign(c *entities.ScheduledCampaign) error {
	err := db.Where("campaign_id = ?", c.CampaignID).Save(c).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return db.Create(c).Error
		}
		return err
	}
	return nil
}

// DeleteScheduledCampaign deletes a scheduled campaign.
func (db *store) DeleteScheduledCampaign(campaignID int64) error {
	return db.Where("campaign_id = ?", campaignID).Delete(entities.ScheduledCampaign{}).Error
}
