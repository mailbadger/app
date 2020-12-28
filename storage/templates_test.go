package storage

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/mailbadger/app/entities"
)

func createTemplates(store Storage) {
	for i := 0; i < 100; i++ {
		err := store.CreateTemplate(&entities.Template{
			Name:        "foo " + strconv.Itoa(i),
			SubjectPart: "Template {{.subject}} " + strconv.Itoa(i),
			UserID:      1,
			TextPart:    "draft {{.text}} " + strconv.Itoa(i),
		})
		if err != nil {
			logrus.Fatal(err)
		}
	}
}

func TestTemplate(t *testing.T) {
	db := openTestDb()
	defer func() {
		err := db.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()
	store := From(db)
	createTemplates(store)

	// templates for insert
	template := &entities.Template{
		UserID:      1,
		Name:        "template1",
		TextPart:    "asd {{.name}}",
		SubjectPart: "subject",
	}

	// test insert templates
	err := store.CreateTemplate(template)
	assert.Nil(t, err)

	// template not found
	templateNotFound, err := store.GetTemplateByName("not-found", 1)
	assert.Equal(t, errors.New("record not found"), err)
	assert.Equal(t, new(entities.Template), templateNotFound)

	template, err = store.GetTemplate(0, 1)
	assert.Equal(t, errors.New("record not found"), err)
	assert.Equal(t, new(entities.Template), template)

	// get template by name and user id test
	templateByName, err := store.GetTemplateByName(template.Name, 1)
	assert.Nil(t, err)
	assert.Equal(t, template.Name, templateByName.Name)
	assert.Equal(t, template.TextPart, templateByName.TextPart)
	assert.Equal(t, template.SubjectPart, templateByName.SubjectPart)

	// get template by id and user id test
	templateByID, err := store.GetTemplate(template.ID, 1)
	assert.Nil(t, err)
	assert.Equal(t, template.Name, templateByID.Name)
	assert.Equal(t, template.TextPart, templateByID.TextPart)
	assert.Equal(t, template.SubjectPart, templateByID.SubjectPart)

	// update template testing
	template.TextPart = "asd {{.name}} and {{.surname}}"
	template.SubjectPart = "Subject {{.update}}"

	err = store.UpdateTemplate(template)
	assert.Nil(t, err)

	templateByID, err = store.GetTemplate(template.ID, 1)
	assert.Nil(t, err)
	assert.Equal(t, template.Name, templateByID.Name)
	assert.Equal(t, template.TextPart, templateByID.TextPart)
	assert.Equal(t, template.SubjectPart, templateByID.SubjectPart)

	p := NewPaginationCursor("/api/templates", 10)
	err = store.GetTemplates(1, p, nil)
	assert.Nil(t, err)
	assert.NotEmpty(t, p.Collection)

	err = store.GetTemplates(1, p, map[string]string{"name": "template"})
	assert.Nil(t, err)
	assert.NotEmpty(t, p.Collection)

	err = store.GetTemplates(1, p, map[string]string{"name": "$%$#"})
	assert.Nil(t, err)
	assert.Empty(t, p.Collection)

	// Test get templates forwards
	p = NewPaginationCursor("/api/templates", 11)
	for i := 0; i < 10; i++ {
		err := store.GetTemplates(1, p, nil)
		assert.Nil(t, err)
		col := p.Collection.(*[]entities.TemplatesCollectionItem)
		assert.NotNil(t, col)
		assert.NotEmpty(t, *col)
		if p.Links.Next != nil {
			assert.Equal(t, 11, len(*col))
			assert.Equal(t, fmt.Sprintf("/api/templates?per_page=%d&starting_after=%d", 11, (*col)[len(*col)-1].GetID()), *p.Links.Next)
			p.SetStartingAfter((*col)[len(*col)-1].GetID())
		} else {
			assert.Equal(t, 2, len(*col))
		}
	}
	assert.Equal(t, int64(101), p.Total)

	// Test get templates backwards
	p = NewPaginationCursor("/api/templates", 13)
	p.SetEndingBefore(1)
	for i := 0; i < 8; i++ {
		err := store.GetTemplates(1, p, nil)
		assert.Nil(t, err)
		col := p.Collection.(*[]entities.TemplatesCollectionItem)
		assert.NotNil(t, col)
		assert.NotEmpty(t, *col)
		if p.Links.Previous != nil {
			assert.Equal(t, 13, len(*col))
			assert.Equal(t, fmt.Sprintf("/api/templates?ending_before=%d&per_page=%d", (*col)[0].GetID(), 13), *p.Links.Previous)
			p.SetEndingBefore((*col)[0].GetID())
		} else {
			assert.Equal(t, 9, len(*col))
		}
	}
	err = store.DeleteTemplate(templateByID.ID, 1)
	assert.Nil(t, err)
}
