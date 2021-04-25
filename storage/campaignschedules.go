package storage

import (
	"errors"
	"time"

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

// GetScheduledCampaigns returns all scheduled campaigns < time
func (db *store) GetScheduledCampaigns(time time.Time) ([]entities.CampaignSchedule, error) {
	var campaignsSchedule []entities.CampaignSchedule
	err := db.Joins("LEFT JOIN campaigns ON campaigns.id = campaign_schedules.campaign_id").
		Where("campaigns.status = ? OR campaigns.status = ? and campaign_schedules.created_at <= ?", entities.StatusDraft, entities.StatusScheduled, time).Find(&campaignsSchedule).Error
	if err != nil {
		return nil, err
	}
	return campaignsSchedule, err
}
