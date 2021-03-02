package storage

import (
	"fmt"

	"github.com/mailbadger/app/entities"
)

// CreateReport creates a report.
func (db *store) CreateCampaignFailedLog(l *entities.CampaignFailedLog) error {
	return db.Create(l).Error
}

func (db *store) LogFailedCampaign(c *entities.Campaign, log *entities.CampaignFailedLog) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Model(&entities.Campaign{}).
		Where("id = ? AND user_id = ?", c.ID, c.UserID).
		Update("status", entities.StatusFailed).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("store: update campaign: %w", err)
	}

	err = tx.Create(log).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("store: create failed campaign log: %w", err)
	}

	return tx.Commit().Error
}
