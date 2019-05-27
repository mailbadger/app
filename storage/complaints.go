package storage

import (
	"github.com/news-maily/app/entities"
)

func (db *store) CreateComplaint(c *entities.Complaint) error {
	return db.Create(c).Error
}
