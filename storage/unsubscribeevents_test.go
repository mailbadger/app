package storage

import (
	"testing"

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

	sub := &entities.Subscriber{
		UserID:      1,
		Name:        "bla",
		Email:       "bla@email.com",
		MetaJSON:    nil,
		Segments:    nil,
		Blacklisted: false,
		Active:      true,
		Errors:      nil,
		Metadata:    nil,
	}

	err := store.CreateSubscriber(sub)
	assert.Nil(t, err)

	err = store.DeactivateSubscriber(sub.UserID, sub.Email)
	assert.Nil(t, err)

	s, err := store.GetSubscriber(sub.UserID, sub.ID)
	assert.Nil(t, err)

	assert.Equal(t, s.Active, false)

}
