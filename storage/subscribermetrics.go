package storage

import (
	"github.com/mailbadger/app/entities"
	"github.com/sirupsen/logrus"
)

func (db *store) UpdateSubscriberMetrics(sm []*entities.SubscribersMetrics, job *entities.Job) (err error) {
	tx := db.DB.Begin()
	
	for _, metric := range sm {
		err = tx.Exec(`INSERT INTO subscribers_metrics
			(user_id, created, deleted, unsubscribed, date)
			VALUES (?, ?, ?, ?, ?)
			ON DUPLICATE KEY
			UPDATE subscribers_metrics
			SET created = ?, deleted = ?, unsubscribed = ?
			WHERE user_id = ? AND date = ?;`,
			metric.UserID, metric.Created, metric.Deleted, metric.Unsubscribed, metric.Date,
			metric.Created, metric.Deleted, metric.Unsubscribed,
			metric.UserID, metric.Date,
		).Error
		if err != nil {
			rbErr := tx.Rollback().Error
			if rbErr != nil {
				logrus.WithError(rbErr).Error("failed to rollback")
			}
		}
		
		return
	}
	
	err = tx.Save(job).Error
	if err != nil {
		rbErr := tx.Rollback().Error
		if rbErr != nil {
			logrus.WithError(rbErr).Error("failed to rollback")
		}
		
		return
	}
	
	return tx.Commit().Error
}
