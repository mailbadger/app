package storage

import (
	"errors"
	"fmt"

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
		Update("status", entities.StatusScheduled).
		Update("event_id", c.ID).Error
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
		Update("status", entities.StatusDraft).
		Update("event_id", nil).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("store: update campaign: %w", err)
	}

	err = tx.Where("campaign_id = ?", campaignID).Delete(entities.CampaignSchedule{CampaignID: campaignID}).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("store: delete campaign schedule: %w", err)
	}

	return tx.Commit().Error
}
