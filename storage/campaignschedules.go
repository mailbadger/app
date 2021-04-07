package storage

import (
	"errors"

	"github.com/jinzhu/gorm"

	"github.com/mailbadger/app/entities"
)

// CreateCampaignSchedule creates a scheduled campaign.
func (db *store) CreateCampaignSchedule(c *entities.CampaignSchedule) error {
	err := db.Where("campaign_id = ?", c.CampaignID).Save(c).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return db.Create(c).Error
		}
		return err
	}
	return nil
}

// DeleteCampaignSchedule deletes a scheduled campaign.
func (db *store) DeleteCampaignSchedule(campaignID int64) error {
	return db.Where("campaign_id = ?", campaignID).Delete(entities.CampaignSchedule{}).Error
}
