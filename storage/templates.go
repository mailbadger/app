package storage

import (
	"github.com/mailbadger/app/entities"
)

// DeleteTemplate deletes the template with given template id and user id from db
func (db *store) DeleteTemplate(templateID int64, userID int64) error {
	return db.Where("user_id = ? and id = ?", userID, templateID).Delete(entities.Template{}).Error
}
