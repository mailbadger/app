package storage

import (
	"fmt"

	"github.com/segmentio/ksuid"

	"github.com/mailbadger/app/entities"
)

func (db *store) LogFailedCampaign(c *entities.Campaign, description string) error {
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

	log := &entities.CampaignFailedLog{
		ID:          ksuid.New(),
		UserID:      c.UserID,
		CampaignID:  c.ID,
		Description: description,
	}

	err = tx.Create(log).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("store: create failed campaign log: %w", err)
	}

	return tx.Commit().Error
}
