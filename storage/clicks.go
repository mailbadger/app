package storage

import "github.com/mailbadger/app/entities"

func (db *store) CreateClick(c *entities.Click) error {
	return db.Create(c).Error
}

// GetCampaignClicksStats fetches collection of clicks stats by campaign id and user id from database
func (db *store) GetCampaignClicksStats(id, userID int64) ([]entities.ClicksStats, error) {
	var clickStats []entities.ClicksStats
	err := db.Table("clicks").
		Select("link, COUNT(DISTINCT(recipient)) AS unique_clicks, COUNT(recipient) AS total_clicks").
		Where("campaign_id = ? AND user_id = ?", id, userID).
		Group("link").
		Find(&clickStats).
		Error

	return clickStats, err
}

// DeleteAllClicksForUser deletes all clicks for user
func (db *store) DeleteAllClicksForUser(userID int64) error {
	return db.Where("user_id = ?", userID).Delete(&entities.Click{}).Error
}

