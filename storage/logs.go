package storage

import "github.com/mailbadger/app/entities"

func (db *store) CreateSendBulkLog(l *entities.SendBulkLog) error {
	return db.Create(l).Error
}

func (db *store) CountLogsByUUID(uuid string) (int, error) {
	var count int
	err := db.Model(&entities.SendBulkLog{}).Where("uuid = ?", uuid).Count(&count).Error
	return count, err
}
