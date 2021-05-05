package storage

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/mailbadger/app/entities"
)

// CreateCampaignSchedule creates a scheduled campaign.
func (db *store) CreateCampaignSchedule(c *entities.CampaignSchedule) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Model(&entities.Campaign{}).
		Where("id = ?", c.CampaignID).
		Update("status", entities.StatusScheduled).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("store: update campaign: %w", err)
	}

	err = tx.Where("campaign_id = ?", c.CampaignID).Save(c).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return fmt.Errorf("store: save campaign schedule: %w", err)
		}
		err = tx.Create(c).Error
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("store: create campaign schedule: %w", err)
		}
	}

	return tx.Commit().Error
}

// DeleteCampaignSchedule deletes a scheduled campaign.
func (db *store) DeleteCampaignSchedule(campaignID int64) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Model(&entities.Campaign{}).
		Where("id = ?", campaignID).
		Update("status", entities.StatusDraft).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("store: update campaign: %w", err)
	}

	err = tx.Where("campaign_id = ?", campaignID).Delete(entities.CampaignSchedule{}).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("store: delete campaign schedule: %w", err)
	}

	return tx.Commit().Error
}

// GetScheduledCampaigns returns all scheduled campaigns < time
func (db *store) GetScheduledCampaigns(time time.Time) ([]entities.CampaignSchedule, error) {
	var campaignsSchedule []entities.CampaignSchedule
	err := db.Joins("JOIN campaigns ON campaigns.id = campaign_schedules.campaign_id").
		Where("campaigns.status = ? and campaign_schedules.scheduled_at <= ?", entities.StatusScheduled, time).Find(&campaignsSchedule).Error
	if err != nil {
		return nil, err
	}
	return campaignsSchedule, err
}
