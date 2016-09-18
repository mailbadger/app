package entities

import (
	"bytes"
	"errors"
	"html/template"
	"time"

	. "github.com/FilipNikolovski/news-maily/middleware"
)

// Template represents the template entity i.e. the email template to be sent
// to subscribers
type Template struct {
	Id        int64     `json:"id" gorm:"column:id; primary_key:yes"`
	UserId    int64     `json:"-" gorm:"column:user_id"`
	Name      string    `json:"name"`
	Content   string    `json:"content" gorm:"column:content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate template properties,
// the template should be able to execute with the given variables
func (t *Template) validate() error {
	switch {
	case t.Name == "":
		return errors.New("Name not specified")
	case t.Content == "":
		return errors.New("Content not specified")
	}

	var buff bytes.Buffer

	td := struct {
		Tracker string
		From    string
	}{
		"<img src='http://example.com/track",
		"John Doe <foo@bar.com>",
	}

	tmpl, err := template.New("html_template").Parse(t.Content)
	if err != nil {
		return err
	}

	return tmpl.Execute(&buff, td)
}

// GetTemplates fetches templates by user id, and populates the pagination obj
func GetTemplates(user_id int64, p *Pagination) {
	var templates []Template
	var count uint64

	db.Offset(p.Offset).Limit(p.PerPage).Where("user_id = ?", user_id).Find(&templates).Count(&count)
	p.SetTotal(count)

	for _, t := range templates {
		p.Append(t)
	}
}

// GetTemplate returns the template by the given id and user id
func GetTemplate(id int64, user_id int64) (Template, error) {
	template := Template{}
	err := db.Where("user_id = ? and id = ?", user_id, id).Find(&template).Error
	return template, err
}

// CreateTemplate
func CreateTemplate(t *Template) error {
	if err := t.validate(); err != nil {
		return err
	}
	return db.Save(t).Error
}

// UpdateTemplate edits an existing template in the database.
func UpdateTemplate(t *Template) error {
	if err := t.validate(); err != nil {
		return err
	}

	return db.Where("id = ?", t.Id).Save(t).Error
}

// DeleteTemplate deletes an existing template in the database.
func DeleteTemplate(id int64, user_id int64) error {
	return db.Where("user_id = ?", user_id).Delete(Template{Id: id}).Error
}
