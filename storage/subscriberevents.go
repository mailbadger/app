package storage

import "github.com/mailbadger/app/entities"

// DeleteAllEventsForUser deletes all subscriber events for user
func (db *store) DeleteAllEventsForUser(userID int64) error {
	return db.Where("user_id = ?", userID).Delete(&entities.SubscriberEvent{}).Error
}
