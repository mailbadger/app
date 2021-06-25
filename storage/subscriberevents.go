package storage

import (
	"github.com/mailbadger/app/entities"
	"github.com/segmentio/ksuid"
)

// DeleteAllEventsForUser deletes all subscriber events for user
func (db *store) DeleteAllEventsForUser(userID int64) error {
	return db.Where("user_id = ?", userID).Delete(&entities.SubscriberEvent{}).Error
}

// GetEventsAfterID fetches limited batch ov events after provided id
func (db *store) GetEventsAfterID(id ksuid.KSUID, limit int64) ([]entities.SubscriberEvent, error) {
	var events []entities.SubscriberEvent
	err := db.Where("id > ?", id).Limit(limit).Find(&events).Error
	return events, err
}
