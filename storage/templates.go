package storage

import (
	"github.com/mailbadger/app/entities"
)

// CreateSubscriber creates a new subscriber in the database.
func (db *store) CreateTemplate(t *entities.Template) error {
	return db.Create(t).Error
}

// UpdateReport edits an existing template in the database.
func (db *store) UpdateTemplate(t *entities.Template) error {
	return db.Where("user_id = ? and id = ?", t.UserID, t.ID).Save(t).Error
}

// GetTemplateByName returns the template by the given name and user id
func (db *store) GetTemplateByName(name string, userID int64) (*entities.Template, error) {
	var template = new(entities.Template)
	err := db.Where("user_id = ? and name = ?", userID, name).First(template).Error
	return template, err
}

// GetTemplate returns the template by the given id and user id
func (db *store) GetTemplate(templateID, userID int64) (*entities.Template, error) {
	var template = new(entities.Template)
	err := db.Where("user_id = ? and id = ?", userID, templateID).First(template).Error
	return template, err
}

// GetTemplates fetches templates by user id, and populates the pagination obj
func (db *store) GetTemplates(userID int64, p *PaginationCursor, scopeMap map[string]string) error {
	p.SetCollection(new([]entities.BaseTemplate))
	p.SetResource("templates")

	p.AddScope(BelongsToUser(userID))
	val, ok := scopeMap["name"]
	if ok {
		p.AddScope(NameLike(val))
	}

	query := db.Table(p.Resource).
		Order("created_at desc, id desc").
		Limit(p.PerPage)

	p.SetQuery(query)

	return db.Paginate(p, userID)
}

// DeleteTemplate deletes the template with given template id and user id from db
func (db *store) DeleteTemplate(templateID int64, userID int64) error {
	return db.Where("user_id = ? and id = ?", userID, templateID).Delete(&entities.Template{}).Error
}

// GetAllTemplatesForUser fetches all templates for user
func (db *store) GetAllTemplatesForUser(userID int64) ([]entities.Template, error) {
	var t []entities.Template
	err := db.Where("user_id = ?", userID).Find(&t).Error
	return t, err
}
