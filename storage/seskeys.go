package storage

import (
	"github.com/news-maily/api/entities"
)

// GetSesKeys returns the SES keys by the given user id
func (db *store) GetSesKeys(userID int64) (*entities.SesKeys, error) {
	var s = new(entities.SesKeys)
	err := db.Where("user_id = ?", userID).First(s).Error
	return s, err
}

// CreateSesKeys adds new SES keys in the database.
func (db *store) CreateSesKeys(s *entities.SesKeys) error {
	return db.Create(s).Error
}
