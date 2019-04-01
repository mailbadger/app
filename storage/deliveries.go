package storage

import (
	"github.com/news-maily/api/entities"
)

func (db *store) CreateDelivery(d *entities.Delivery) error {
	return db.Create(d).Error
}
