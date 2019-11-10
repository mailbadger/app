package storage

import (
	"testing"
	"time"

	"github.com/news-maily/app/entities"
	"github.com/news-maily/app/utils/pagination"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestSubscriber(t *testing.T) {
	db := openTestDb()
	defer func() {
		err := db.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()
	store := From(db)

	l := &entities.Segment{
		Name:   "foo",
		UserID: 1,
	}

	err := store.CreateSegment(l)
	assert.Nil(t, err)

	//Test create subscriber
	s := &entities.Subscriber{
		Name:        "foo",
		Email:       "john@example.com",
		UserID:      1,
		MetaJSON:    []byte(`{"foo":"bar"}`),
		Blacklisted: false,
		Active:      true,
	}
	s.Segments = append(s.Segments, *l)

	err = store.CreateSubscriber(s)
	assert.Nil(t, err)

	//Test get subscriber
	s, err = store.GetSubscriber(s.ID, 1)
	s.Normalize()

	assert.Nil(t, err)
	assert.Equal(t, s.Name, "foo")
	assert.NotEmpty(t, s.Metadata)
	assert.Equal(t, s.Metadata["foo"], "bar")

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
	assert.Equal(t, s.Errors["email"], "The specified email is not valid.")

	//Test get subs
	cp := &pagination.Cursor{PerPage: 10, StartingAfter: 0}
	store.GetSubscribers(1, cp)
	assert.NotEmpty(t, cp.Collection)

	//Test get subs by ids
	subs, err := store.GetSubscribersByIDs([]int64{1}, 1)
	assert.Nil(t, err)
	assert.NotEmpty(t, subs)

	//Test get subs by list id
	p := &pagination.Cursor{PerPage: 10}
	store.GetSubscribersBySegmentID(l.ID, 1, p)
	assert.NotEmpty(t, p.Collection)

	var timestamp time.Time
	subs, err = store.GetDistinctSubscribersBySegmentIDs([]int64{l.ID}, 1, false, true, timestamp, 0, 10)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(subs))

	err = store.DeleteSubscriber(1, 1)
	assert.Nil(t, err)
}
