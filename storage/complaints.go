package storage

import (
	"github.com/mailbadger/app/entities"
)

func (db *store) CreateComplaint(c *entities.Complaint) error {
	return db.Create(c).Error
}

// DeleteAllComplaintsForUser deletes all complaints for user
func (db *store) DeleteAllComplaintsForUser(userID int64) error {
	return db.Where("user_id = ?", userID).Delete(&entities.Complaint{}).Error
}
