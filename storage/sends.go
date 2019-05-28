package storage

import "github.com/news-maily/app/entities"

func (db *store) CreateSend(s *entities.Send) error {
	return db.Create(s).Error
}
