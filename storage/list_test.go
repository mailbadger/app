package storage

import (
	"testing"

	"github.com/FilipNikolovski/news-maily/entities"
	"github.com/FilipNikolovski/news-maily/utils/pagination"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	db := openTestDb()
	defer db.Close()

	store := From(db)

	//Test create list
	l := &entities.List{
		Name:   "foo",
		UserId: 1,
		Metadata: []entities.ListMetadata{
			{Key: "key1", Value: "val1"},
		},
	}

	err := store.CreateList(l)
	assert.Nil(t, err)

	//Test get list
	l, err = store.GetList(1, 1)
	assert.Nil(t, err)
	assert.Equal(t, l.Name, "foo")
	assert.NotEmpty(t, l.Metadata)
	assert.Equal(t, l.Metadata[0].Key, "key1")
	assert.Equal(t, l.Metadata[0].Value, "val1")

	//Test update list
	l.Name = "bar"
	err = store.UpdateList(l)
	assert.Nil(t, err)
	assert.Equal(t, l.Name, "bar")

	//Test list validation when name and subject are invalid
	l.Name = ""
	l.Validate()
	assert.Equal(t, l.Errors["name"], entities.ErrListNameEmpty.Error())

	//Test get lists
	p := &pagination.Pagination{}
	store.GetLists(1, p)
	assert.NotEmpty(t, p.Collection)
	assert.Equal(t, len(p.Collection), int(p.Total))

	// Test delete list
	err = store.DeleteList(1, 1)
	assert.Nil(t, err)
}
