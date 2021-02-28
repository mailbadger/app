package storage

import (
	"testing"
	"time"

	"github.com/segmentio/ksuid"
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
	id := ksuid.New()

	unsubscribeEvents := []*entities.UnsubscribeEvent{
		{
			Email:     "email1@bla.com",
			UserID:    1,
			CreatedAt: now,
		},
		{
			Email:     "email2@bla.com",
			UserID:    1,
			CreatedAt: now,
		},
		{
			Email:     "email3@bla.com",
			UserID:    1,
			CreatedAt: now,
		},
	}
	// test insert opens
	for i, k := range unsubscribeEvents {
		k.ID = id
		err := store.CreateUnsubscribeEvent(unsubscribeEvents[i])
		assert.Nil(t, err)
		id = id.Next()
	}

}
