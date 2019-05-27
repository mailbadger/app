package storage

import "github.com/news-maily/app/entities"

func (db *store) CreateClick(c *entities.Click) error {
	return db.Create(c).Error
}
