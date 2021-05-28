package storage

import (
	"github.com/mailbadger/app/entities"
)

func (db *store) CreateSendLog(l *entities.SendLog) error {
	return db.Create(l).Error
}

func (db *store) CountLogsByUUID(id string) (int, error) {
	var count int
	err := db.Model(&entities.SendLog{}).Where("id = ?", id).Count(&count).Error
	return count, err
}

func (db *store) CountLogsByStatus(status string) (int, error) {
	var count int
	err := db.Model(&entities.SendLog{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

// GetSendLogByUUID returns send log with specified uuid
func (db *store) GetSendLogByUUID(id string) (*entities.SendLog, error) {
	var log = new(entities.SendLog)
	err := db.Where("id = ?", id).Find(log).Error
	return log, err
}

// DeleteAllSendLogsForUser deletes all send log records for user
func (db *store) DeleteAllSendLogsForUser(userID int64) error {
	return db.Where("user_id = ?", userID).Delete(&entities.SendLog{}).Error
}
