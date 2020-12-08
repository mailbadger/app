package storage

import (
	"github.com/mailbadger/app/entities"
)

func (db *store) GetTemplate(templateID int64, userID int64) (*entities.Template, error) {
	var template = new(entities.Template)
	err := db.Where("user_id = ? and id = ?", userID, templateID).Find(template).Error
	return template, err
}
