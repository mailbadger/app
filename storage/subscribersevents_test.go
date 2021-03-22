package storage

import (
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/mailbadger/app/entities"
)

func TestSubscribersEvents(t *testing.T) {
	db := openTestDb()
	defer func() {
		err := db.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()
	store := From(db)

	// create subscriber
	s := &entities.Subscriber{
		Name:        "foo",
		Email:       "john@example.com",
		UserID:      1,
		MetaJSON:    []byte(`{"foo":"bar"}`),
		Blacklisted: false,
		Active:      true,
	}

	err := store.CreateSubscriber(s)
	assert.Nil(t, err)

	se := &entities.SubscribersEvent{
		ID:              ksuid.New(),
		UserID:          1,
		SubscriberID:    s.ID,
		SubscriberEmail: s.Email,
		EventType:       entities.SubscriberEventTypeCreated,
	}

	err = store.CreateSubscribersEvent(se)
	assert.Nil(t, err)
}
