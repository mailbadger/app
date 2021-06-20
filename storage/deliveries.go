package storage

import (
	"github.com/mailbadger/app/entities"
)

func (db *store) CreateDelivery(d *entities.Delivery) error {
	return db.Create(d).Error
}

// DeleteAllDeliveriesForUser deletes all deliveries for user
func (db *store) DeleteAllDeliveriesForUser(userID int64) error {
	return db.Where("user_id = ?", userID).Delete(&entities.Delivery{}).Error
}
