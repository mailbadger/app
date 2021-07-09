package storage

import (
	"testing"
	
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	
	"github.com/mailbadger/app/entities"
)

func TestSubscriberEvents(t *testing.T) {
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
}
