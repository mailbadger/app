package storage

import (
	"errors"
	"testing"

	"github.com/FilipNikolovski/news-maily/entities"
	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	db := openTestDb()
	defer db.Close()

	store := From(db)

	//Test create template
	err := store.CreateTemplate(&entities.Template{
		Name:    "foo",
		UserId:  1,
		Content: "Foo bar",
	})

	assert.Nil(t, err)

	//Test get template
	template, err := store.GetTemplate(1, 1)
	assert.Nil(t, err)
	assert.Equal(t, template.Name, "foo")
	assert.Equal(t, template.Content, "Foo bar")

	//Test create template when name is empty
	err = store.CreateTemplate(&entities.Template{
		UserId:  1,
		Content: "Foo bar",
	})

	assert.Equal(t, err, errors.New("Name not specified"))

	//Test create template when content is empty
	err = store.CreateTemplate(&entities.Template{
		Name:   "foo",
		UserId: 1,
	})

	assert.Equal(t, err, errors.New("Content not specified"))
}
