package storage

import (
	"testing"

	"github.com/FilipNikolovski/news-maily/entities"
	"github.com/FilipNikolovski/news-maily/utils/pagination"
	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	db := openTestDb()
	defer db.Close()

	store := From(db)

	//Test create template
	template := &entities.Template{
		Name:    "foo",
		UserId:  1,
		Content: "Foo bar",
	}

	err := store.CreateTemplate(template)

	assert.Nil(t, err)

	//Test get template
	template, err = store.GetTemplate(1, 1)
	assert.Nil(t, err)
	assert.Equal(t, template.Name, "foo")
	assert.Equal(t, template.Content, "Foo bar")

	//Test get template by name
	template, err = store.GetTemplateByName("foo", 1)
	assert.Nil(t, err)
	assert.Equal(t, template.Name, "foo")
	assert.Equal(t, template.Content, "Foo bar")

	//Test update template
	template.Name = "bar"
	err = store.UpdateTemplate(template)
	assert.Nil(t, err)
	assert.Equal(t, template.Name, "bar")

	//Test template validation when name is invalid
	template.Name = ""
	template.Validate()
	assert.Equal(t, template.Errors["name"], entities.ErrTemplateNameEmpty.Error())

	//Test get templates
	p := &pagination.Pagination{}
	store.GetTemplates(1, p)
	assert.NotEmpty(t, p.Collection)
	assert.Equal(t, len(p.Collection), int(p.Total))

	// Test delete template
	err = store.DeleteTemplate(1, 1)
	assert.Nil(t, err)
}
