package storage

import "github.com/news-maily/api/entities"

func (db *store) CreateOpen(o *entities.Open) error {
	return db.Create(o).Error
}
