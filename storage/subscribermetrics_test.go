package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mailbadger/app/entities"
)

func TestSubscriberMetrics(t *testing.T) {
	db := openTestDb()

	store := From(db)

	sm := &entities.SubscriberMetrics{
		UserID:       1,
		Created:      13,
		Unsubscribed: 6,
		Date:         time.Now(),
	}

	err := store.UpdateSubscriberMetrics(sm)
	assert.Nil(t, err)
	sm.Created = 23
	sm.Unsubscribed = 23

	err = store.UpdateSubscriberMetrics(sm)
	assert.Nil(t, err)
}
