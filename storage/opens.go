package storage

import "github.com/mailbadger/app/entities"

func (db *store) CreateOpen(o *entities.Open) error {
	return db.Create(o).Error
}
