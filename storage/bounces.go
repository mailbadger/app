package storage

import (
	"github.com/news-maily/app/entities"
)

func (db *store) CreateBounce(b *entities.Bounce) error {
	return db.Create(b).Error
}
