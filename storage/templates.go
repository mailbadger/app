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
	err := db.Where("user_id = ? and name = ?", userID, name).Find(template).Error
	return template, err
}

// GetTemplate returns the template by the given id and user id
func (db *store) GetTemplate(templateID, userID int64) (*entities.Template, error) {
	var template = new(entities.Template)
	err := db.Where("user_id = ? and id = ?", userID, templateID).Find(template).Error
	return template, err
}

// GetTemplates fetches templates by user id, and populates the pagination obj
func (db *store) GetTemplates(userID int64, p *PaginationCursor, scopeMap map[string]string) error {
	p.SetCollection(&[]entities.TemplatesCollectionItem{})
	p.SetResource("templates")

	for k, v := range scopeMap {
		if k == "name" {
			p.AddScope(NameLike(v))
		}
	}

	query := db.Table(p.Resource).
		Where("user_id = ?", userID).
		Order("created_at desc, id desc").
		Limit(p.PerPage)

	p.SetQuery(query)

	return db.Paginate(p, userID)
}

// DeleteTemplate deletes the template with given template id and user id from db
func (db *store) DeleteTemplate(templateID int64, userID int64) error {
	return db.Delete(entities.Template{Model: entities.Model{ID: templateID}, UserID: userID}).Error
}
