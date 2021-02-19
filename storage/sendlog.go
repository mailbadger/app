package storage

import (
	"github.com/mailbadger/app/entities"
)

func (db *store) CreateSendLog(l *entities.SendLog) error {
	return db.Create(l).Error
}

func (db *store) CountLogsByUUID(uid string) (int, error) {
	var count int
	err := db.Model(&entities.SendLog{}).Where("uid = ?", uid).Count(&count).Error
	return count, err
}

func (db *store) CountLogsByStatus(status string) (int, error) {
	var count int
	err := db.Model(&entities.SendLog{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

// GetSendLogByUUID returns send log with specified uuid
func (db *store) GetSendLogByUUID(uid string) (*entities.SendLog, error) {
	var log = new(entities.SendLog)
	err := db.Where("uid = ?", uid).Find(log).Error
	return log, err
}
