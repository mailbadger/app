package storage

import (
	"time"
	
	"github.com/mailbadger/app/entities"
)

// DeleteAllEventsForUser deletes all subscriber events for user
func (db *store) DeleteAllEventsForUser(userID int64) error {
	return db.Where("user_id = ?", userID).Delete(&entities.SubscriberEvent{}).Error
}

// GetGroupedSubscriberEvents fetches batch ov events between two dates
func (db *store) GetGroupedSubscriberEvents(startDate time.Time, endDate time.Time) ([]*entities.GroupedSubscriberEvents, error) {
	var events []*entities.GroupedSubscriberEvents
	err := db.Raw(`SELECT user_id, DATE(created_at) date, event_type, COUNT(id) total
		FROM subscriber_events
		WHERE created_at >= ? AND created_at <= ?
		GROUP BY user_id, DATE(created_at), event_type;`,
		startDate, endDate).Scan(&events).Error
	return events, err
}
