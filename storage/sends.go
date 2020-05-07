package storage

import "github.com/mailbadger/app/entities"

func (db *store) CreateSend(s *entities.Send) error {
	return db.Create(s).Error
}
