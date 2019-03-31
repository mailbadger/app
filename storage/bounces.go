package storage

import (
	"github.com/news-maily/api/entities"
)

func (db *store) CreateBounce(b *entities.Bounce) error {
	return db.Create(b).Error
}
