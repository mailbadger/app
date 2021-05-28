package storage

import (
	"github.com/mailbadger/app/entities"
)

func (db *store) CreateBounce(b *entities.Bounce) error {
	return db.Create(b).Error
}

// DeleteAllBouncesForUser deletes all bounces for user
func (db *store) DeleteAllBouncesForUser(userID int64) error {
	return db.Where("user_id = ?", userID).Delete(&entities.Bounce{}).Error
}
