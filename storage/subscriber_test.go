package storage

import (
	"testing"

	"github.com/news-maily/api/entities"
	"github.com/news-maily/api/utils/pagination"
	"github.com/stretchr/testify/assert"
)

func TestSubscriber(t *testing.T) {
	db := openTestDb()
	defer db.Close()

	store := From(db)

	//Test create subscriber
	s := &entities.Subscriber{
		Name:   "foo",
		Email:  "john@example.com",
		UserId: 1,
		Metadata: []entities.SubscriberMetadata{
			{Key: "key", Value: "val"},
		},
	}

	err := store.CreateSubscriber(s)
	assert.Nil(t, err)

	//Test get subscriber
	s, err = store.GetSubscriber(1, 1)
	assert.Nil(t, err)
	assert.Equal(t, s.Name, "foo")
	assert.NotEmpty(t, s.Metadata)
	assert.Equal(t, s.Metadata[0].Key, "key")
	assert.Equal(t, s.Metadata[0].Value, "val")

	//Test get subscriber by email
	s, err = store.GetSubscriberByEmail("john@example.com", 1)
	assert.Nil(t, err)
	assert.Equal(t, s.Name, "foo")

	//Test update subscriber
	s.Name = "bar"
	err = store.UpdateSubscriber(s)
	assert.Nil(t, err)
	assert.Equal(t, s.Name, "bar")

	//Test subscriber validation when name and email are invalid
	s.Name = ""
	s.Email = "foo bar"
	s.Validate()
	assert.Equal(t, s.Errors["name"], entities.ErrSubscriberNameEmpty.Error())
	assert.Equal(t, s.Errors["email"], entities.ErrEmailInvalid.Error())

	//Test get subs
	p := &pagination.Pagination{PerPage: 10}
	store.GetSubscribers(1, p)
	assert.NotEmpty(t, p.Collection)
	assert.Equal(t, len(p.Collection), int(p.Total))

	// Test delete subscriber
	err = store.DeleteSubscriber(1, 1)
	assert.Nil(t, err)
}
