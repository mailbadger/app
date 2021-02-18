package storage

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"

	"github.com/mailbadger/app/entities"
)

func TestUnsubscribedSubscriber(t *testing.T) {
	db := openTestDb()
	defer func() {
		err := db.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	store := From(db)
	now := time.Now().UTC()

	unsubscribedSubscribers := []entities.UnsubscribedSubscriber{
		{
			ID:        1,
			Email:     "email1@bla.com",
			CreatedAt: now,
		},
		{
			ID:        2,
			Email:     "email2@bla.com",
			CreatedAt: now,
		},
		{
			ID:        3,
			Email:     "email3@bla.com",
			CreatedAt: now,
		},
	}
	// test insert opens
	for i := range unsubscribedSubscribers {
		err := store.CreateUnsubscribedSubscriber(&unsubscribedSubscribers[i])
		assert.Nil(t, err)
	}

}
