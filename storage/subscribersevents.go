package storage

import "github.com/mailbadger/app/entities"

// CreateSubscribersEvent adds new subscribers event record in database
func (db *store) CreateSubscribersEvent(se *entities.SubscribersEvent) error {
	return db.Create(se).Error
}
