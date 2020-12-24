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
func (db *store) GetTemplate(id, userID int64) (*entities.Template, error) {
	var template = new(entities.Template)
	err := db.Where("user_id = ? and id = ?", userID, id).Find(template).Error
	return template, err
}

// ListTemplates fetches templates by user id, and populates the pagination obj
func (db *store) ListTemplates(userID int64, p *PaginationCursor, scopeMap map[string]string) error {
	p.SetCollection(&[]entities.TemplatesCollection{})
	p.SetResource("templates")

	for k, v := range scopeMap {
		if k == "name" {
			p.AddScope(NameLike(v))
		}
	}

	p.SetQuery(db.Table(p.Resource).
		Where("user_id = ?", userID).
		Order("created_at desc, id desc").
		Limit(p.PerPage))

	return db.Paginate(p, userID)
}
