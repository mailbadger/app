package storage

import (
	"github.com/mailbadger/app/entities"
)

func (db *store) CreateDelivery(d *entities.Delivery) error {
	return db.Create(d).Error
}
