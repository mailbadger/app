package storage

import (
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/mailbadger/app/entities"
)

func TestTemplate(t *testing.T) {
	db := openTestDb()
	defer func() {
		err := db.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()
	store := From(db)

	//templates for insert
	templates := []entities.Template{
		{
			Model :entities.Model{
				ID:        1,
			},
			UserID:      1,
			Name:        "template1",
			TextPart:    "asd {{.name}}",
			SubjectPart: "subject",
		},
		{
			Model :entities.Model{
				ID:        2,
			},
			UserID:      1,
			Name:        "template2",
			TextPart:    "asd {{.name}}",
			SubjectPart: "subject2",
		},
	}

	// test insert templates
	for _, te := range templates {
		err := store.CreateTemplate(&te)
		assert.Nil(t, err)
	}

	// template not found
	template, err := store.GetTemplateByName("not-found", 1)
	assert.Equal(t, errors.New("record not found"), err)
	assert.Equal(t, new(entities.Template), template)

	// get template by name and user id test
	template, err = store.GetTemplateByName(templates[0].Name, 1)
	assert.Nil(t, err)
	assert.Equal(t, templates[0].Name, template.Name)
	assert.Equal(t, templates[0].TextPart, template.TextPart)
	assert.Equal(t, templates[0].SubjectPart, template.SubjectPart)

	// get template by id and user id test
	template, err = store.GetTemplate(templates[1].ID, 1)
	assert.Nil(t, err)
	assert.Equal(t, templates[1].Name, template.Name)
	assert.Equal(t, templates[1].TextPart, template.TextPart)
	assert.Equal(t, templates[1].SubjectPart, template.SubjectPart)

	templates[1] = entities.Template{
		UserID:      1,
		Name:        "template2",
		TextPart:    "asd {{.name}} and {{.surname}}",
		SubjectPart: "subject2",
	}

	err = store.UpdateTemplate(&templates[1])
	assert.Nil(t, err)

	p := NewPaginationCursor("/api/templates", 10)
	err = store.ListTemplates(1, p, nil)
	assert.Nil(t, err)
	assert.NotEmpty(t, p.Collection)

	err = store.ListTemplates(1, p, map[string]string{"name":"template"})
	assert.Nil(t, err)
	assert.NotEmpty(t, p.Collection)

	err = store.ListTemplates(1, p, map[string]string{"name":"$%$#"})
	assert.Nil(t, err)
	assert.Empty(t, p.Collection)
}
