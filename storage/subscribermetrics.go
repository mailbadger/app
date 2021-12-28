package storage

import (
	"github.com/mailbadger/app/entities"
	"gorm.io/gorm/clause"
)

func (db *store) UpdateSubscriberMetrics(sm *entities.SubscriberMetrics) (err error) {
	return db.Debug().
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"created", "unsubscribed"}),
		}).
		Create(sm).Error
}
