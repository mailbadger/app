package storage

import (
	"github.com/FilipNikolovski/news-maily/entities"
	"github.com/FilipNikolovski/news-maily/utils/pagination"
)

// GetTemplates fetches templates by user id, and populates the pagination obj
func (db *store) GetTemplates(userID int64, p *pagination.Pagination) {
	var templates []entities.Template
	var count uint64

	db.Offset(p.Offset).Limit(p.PerPage).Where("user_id = ?", userID).Find(&templates).Count(&count)
	p.SetTotal(count)

	for _, t := range templates {
		p.Append(t)
	}
}

// GetTemplate returns the template by the given id and user id
func (db *store) GetTemplate(id, userID int64) (*entities.Template, error) {
	var template = new(entities.Template)
	err := db.Where("user_id = ? and id = ?", userID, id).Find(template).Error
	return template, err
}

// GetTemplateByName returns the template by the given name and user id
func (db *store) GetTemplateByName(name string, userID int64) (*entities.Template, error) {
	var template = new(entities.Template)
	err := db.Where("user_id = ? and name = ?", userID, name).Find(template).Error
	return template, err
}

// CreateTemplate creates a new template in the database.
func (db *store) CreateTemplate(t *entities.Template) error {
	return db.Create(t).Error
}

// UpdateTemplate edits an existing template in the database.
func (db *store) UpdateTemplate(t *entities.Template) error {
	return db.Where("id = ? and user_id = ?", t.Id, t.UserId).Save(t).Error
}

// DeleteTemplate deletes an existing template from the database.
func (db *store) DeleteTemplate(id, userID int64) error {
	return db.Where("user_id = ?", userID).Delete(entities.Template{Id: id}).Error
}
