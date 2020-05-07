package storage

import (
	"github.com/mailbadger/app/entities"
)

// GetSesKeys returns the SES keys by the given user id
func (db *store) GetSesKeys(userID int64) (*entities.SesKeys, error) {
	var s = new(entities.SesKeys)
	err := db.Where("user_id = ?", userID).First(s).Error
	if err != nil {
		return nil, err
	}
	return s, nil
}

// CreateSesKeys adds new SES keys in the database.
func (db *store) CreateSesKeys(s *entities.SesKeys) error {
	return db.Create(s).Error
}

// DeleteSesKeys deletes the keys by the given user id.
func (db *store) DeleteSesKeys(userID int64) error {
	return db.Delete(&entities.SesKeys{UserID: userID}).Error
}
