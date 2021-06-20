package storage

import "github.com/mailbadger/app/entities"

func (db *store) CreateOpen(o *entities.Open) error {
	return db.Create(o).Error
}

// DeleteAllOpensForUser deletes all opens for user
func (db *store) DeleteAllOpensForUser(userID int64) error {
	return db.Where("user_id = ?", userID).Delete(&entities.Open{}).Error
}
