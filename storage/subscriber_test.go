package storage

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/mailbadger/app/entities"
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

	s2 := &entities.Subscriber{
		Name:        "foo 2",
		Email:       "john+1@example.com",
		UserID:      1,
		MetaJSON:    []byte(`{"foo":"bar"}`),
		Blacklisted: false,
		Active:      true,
	}

	err = store.CreateSubscriber(s2)
	assert.Nil(t, err)

	//Test get subscriber
	s, err = store.GetSubscriber(s.ID, 1)
	assert.Nil(t, err)

	assert.Equal(t, s.Name, "foo")
	assert.NotEmpty(t, s.MetaJSON)

	var m map[string]string
	err = json.Unmarshal(s.MetaJSON, &m)
	assert.Nil(t, err)
	assert.Equal(t, m["foo"], "bar")

	//Test get subs by list id
	p := NewPaginationCursor(fmt.Sprintf("/api/segments/%d/subscribers", l.ID), 10)
	err = store.GetSubscribersBySegmentID(l.ID, 1, p)
	assert.Nil(t, err)
	assert.NotEmpty(t, p.Collection)

	var timestamp time.Time
	subs, err := store.GetDistinctSubscribersBySegmentIDs([]int64{l.ID}, 1, false, true, timestamp, 0, 10)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(subs))

	//Test get total subs in segment
	totalInSeg, err := store.GetTotalSubscribersBySegment(l.ID, 1)
	assert.Nil(t, err)
	assert.Equal(t, totalInSeg, int64(1))

	//Test get subscriber by email
	s, err = store.GetSubscriberByEmail("john@example.com", 1)
	assert.Nil(t, err)
	assert.Equal(t, s.Name, "foo")

	//Test update subscriber
	s.Name = "bar"
	s.MetaJSON = []byte(`{"foo": "baz"}`)
	err = store.UpdateSubscriber(s)
	assert.Nil(t, err)
	assert.Equal(t, s.Name, "bar")

	m, err = s.GetMetadata()
	assert.Nil(t, err)
	assert.Equal(t, m["foo"], "baz")

	//Test subscriber validation when name and email are invalid
	s.Name = ""
	s.Email = "foo bar"
	s.Validate()
	assert.Equal(t, s.Errors["email"], "The specified email is not valid.")

	//Test get subs
	p = NewPaginationCursor("/api/subcribers", 10)
	err = store.GetSubscribers(1, p)
	assert.Nil(t, err)
	assert.NotEmpty(t, p.Collection)

	//Test get subs by ids
	subs, err = store.GetSubscribersByIDs([]int64{1}, 1)
	assert.Nil(t, err)
	assert.NotEmpty(t, subs)

	//Test get total subs
	total, err := store.GetTotalSubscribers(1)
	assert.Nil(t, err)
	assert.Equal(t, total, int64(2))

	//Test delete subscriber
	err = store.DeleteSubscriberByEmail(s2.Email, 1)
	assert.Nil(t, err)

	err = store.DeleteSubscriber(1, 1)
	assert.Nil(t, err)
}
