package storage

import "github.com/mailbadger/app/entities"

func (db *store) CreateSendLog(l *entities.SendLog) error {
	return db.Create(l).Error
}

func (db *store) CountLogsByUUID(uuid string) (int, error) {
	var count int
	err := db.Model(&entities.SendLog{}).Where("uuid = ?", uuid).Count(&count).Error
	return count, err
}

func (db *store) CountLogsByStatus(status string) (int, error) {
	var count int
	err := db.Model(&entities.SendLog{}).Where("status = ?", status).Count(&count).Error
	return count, err
}
