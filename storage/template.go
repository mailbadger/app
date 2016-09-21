package storage

import (
	"github.com/FilipNikolovski/news-maily/entities"
	"github.com/FilipNikolovski/news-maily/routes/middleware"
)

// GetTemplates fetches templates by user id, and populates the pagination obj
func (db *store) GetTemplates(user_id int64, p *middleware.Pagination) {
	var templates []entities.Template
	var count uint64

	db.Offset(p.Offset).Limit(p.PerPage).Where("user_id = ?", user_id).Find(&templates).Count(&count)
	p.SetTotal(count)

	for _, t := range templates {
		p.Append(t)
	}
}

// GetTemplate returns the template by the given id and user id
func (db *store) GetTemplate(id int64, user_id int64) (entities.Template, error) {
	template := entities.Template{}
	err := db.Where("user_id = ? and id = ?", user_id, id).Find(&template).Error
	return template, err
}

// CreateTemplate
func (db *store) CreateTemplate(t *entities.Template) error {
	if err := t.Validate(); err != nil {
		return err
	}
	return db.Save(t).Error
}

// UpdateTemplate edits an existing template in the database.
func (db *store) UpdateTemplate(t *entities.Template) error {
	if err := t.Validate(); err != nil {
		return err
	}

	return db.Where("id = ?", t.Id).Save(t).Error
}

// DeleteTemplate deletes an existing template in the database.
func (db *store) DeleteTemplate(id int64, user_id int64) error {
	return db.Where("user_id = ?", user_id).Delete(entities.Template{Id: id}).Error
}
