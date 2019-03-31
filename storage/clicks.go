package storage

import "github.com/news-maily/api/entities"

func (db *store) CreateClick(c *entities.Click) error {
	return db.Create(c).Error
}
