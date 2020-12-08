package storage

import (
	"github.com/mailbadger/app/entities"
)

// CreateSubscriber creates a new subscriber in the database.
func (db *store) CreateTemplate(t *entities.Template) error {
	return db.Create(t).Error
}

// UpdateReport edits an existing template in the database.
func (db *store) UpdateTemplate(t *entities.Template) error {
	return db.Where("user_id = ? and name = ?", t.UserID, t.Name).Save(t).Error
}
