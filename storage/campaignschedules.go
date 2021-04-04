package storage

import (
	"errors"

	"github.com/jinzhu/gorm"

	"github.com/mailbadger/app/entities"
)

// CreateScheduledCampaign creates a scheduled campaign.
func (db *store) CreateScheduledCampaign(c *entities.CampaignSchedule) error {
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
	return db.Where("campaign_id = ?", campaignID).Delete(entities.CampaignSchedule{}).Error
}

// GetScheduledCampaign fetches the schedule record for campaign
func (db *store) GetScheduledCampaign(campaignID int64) (*entities.CampaignSchedule, error) {
	var sc = new(entities.CampaignSchedule)
	err := db.Where("campaign_id = ?", campaignID).Find(&sc).Error
	return sc, err
}
