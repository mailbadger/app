package storage

import "github.com/mailbadger/app/entities"

func (db *store) CreateSend(s *entities.Send) error {
	return db.Create(s).Error
}

// DeleteAllSendsForUser deletes all sends for user
func (db *store) DeleteAllSendsForUser(userID int64) error {
	return db.Where("user_id = ?", userID).Delete(&entities.Send{}).Error
}
