package storage

import (
	"github.com/news-maily/app/entities"
)

// GetAccessKeys fetches access keys by user id
func (db *store) GetAccessKeys(userID int64) []*entities.AccessKey {
	var keys []*entities.AccessKey

	db.Where("user_id = ?", userID).Find(&keys)

	return keys
}

// GetAccessKeys fetches access keys by user id
func (db *store) GetAccessKey(identifier string) (*entities.AccessKey, error) {
	var key = new(entities.AccessKey)
	err := db.
		Where("access_key = ? and active = ?", identifier, true).
		Preload("User").
		Find(key).
		Error

	return key, err
}

// CreateAccessKey creates a new access key in the database.
func (db *store) CreateAccessKey(ak *entities.AccessKey) error {
	return db.Create(ak).Error
}

// UpdateAccessKey edits an existing access key in the database.
func (db *store) UpdateAccessKey(ak *entities.AccessKey) error {
	return db.Where("id = ? and user_id = ?", ak.ID, ak.UserID).Save(ak).Error
}

// DeleteAccessKey deletes an existing access key from the database.
func (db *store) DeleteAccessKey(id, userID int64) error {
	return db.Where("user_id = ?", userID).Delete(entities.AccessKey{ID: id}).Error
}
