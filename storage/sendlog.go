package storage

import (
	"github.com/mailbadger/app/entities"
)

func (db *store) CreateSendLog(l *entities.SendLog) error {
	return db.Create(l).Error
}

func (db *store) CountLogsByUUID(id string) (int64, error) {
	var count int64
	err := db.Model(&entities.SendLog{}).Where("id = ?", id).Count(&count).Error
	return count, err
}

func (db *store) CountLogsByStatus(status string) (int64, error) {
	var count int64
	err := db.Model(&entities.SendLog{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

// GetSendLogByUUID returns send log with specified uuid
func (db *store) GetSendLogByUUID(id string) (*entities.SendLog, error) {
	var log = new(entities.SendLog)
	err := db.Where("id = ?", id).First(log).Error
	return log, err
}
