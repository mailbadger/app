package storage

import "github.com/news-maily/api/entities"

func (db *store) CreateSendBulkLog(l *entities.SendBulkLog) error {
	return db.Create(l).Error
}
