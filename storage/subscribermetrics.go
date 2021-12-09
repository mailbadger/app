package storage

import (
	"github.com/mailbadger/app/entities"
)

func (db *store) UpdateSubscriberMetrics(sm *entities.SubscriberMetrics) (err error) {
	return db.Exec(`INSERT INTO subscriber_metrics
			(user_id, created, deleted, unsubscribed, date)
			VALUES (?, ?, ?, ?, ?)
			ON DUPLICATE KEY
			UPDATE created = ?, deleted = ?, unsubscribed = ?`,
		sm.UserID, sm.Created, sm.Deleted, sm.Unsubscribed, sm.Date,
		sm.Created, sm.Deleted, sm.Unsubscribed,
	).Error
}
