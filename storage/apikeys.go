package storage

import (
	"github.com/mailbadger/app/entities"
)

// GetAPIKeys fetches api keys by user id.
func (db *store) GetAPIKeys(userID int64) ([]*entities.APIKey, error) {
	var keys []*entities.APIKey
	err := db.Where("user_id = ?", userID).Find(&keys).Error
	return keys, err
}

// GetAPIKey fetches access keys by the given secret.
func (db *store) GetAPIKey(identifier string) (*entities.APIKey, error) {
	var key = new(entities.APIKey)
	err := db.
		Where("secret_key = ? and active = ?", identifier, true).
		Preload("User.Boundaries").Preload("User.Roles").
		First(key).
		Error

	return key, err
}

// CreateAPIKey creates a new api key in the database.
func (db *store) CreateAPIKey(ak *entities.APIKey) error {
	return db.Create(ak).Error
}

// UpdateAccessKey edits an existing api key in the database.
func (db *store) UpdateAPIKey(ak *entities.APIKey) error {
	return db.Where("id = ? and user_id = ?", ak.ID, ak.UserID).Save(ak).Error
}

// DeleteAccessKey deletes an existing api key from the database.
func (db *store) DeleteAPIKey(id, userID int64) error {
	return db.Where("user_id = ?", userID).Delete(entities.APIKey{ID: id}).Error
}
