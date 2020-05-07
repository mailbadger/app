package storage

import (
	"github.com/mailbadger/app/entities"
)

func (db *store) CreateBounce(b *entities.Bounce) error {
	return db.Create(b).Error
}
