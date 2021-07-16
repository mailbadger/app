package storage

import (
	"testing"
	"time"
	
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	
	"github.com/mailbadger/app/entities"
)

func TestSubscriberMetrics(t *testing.T) {
	db := openTestDb()
	defer func() {
		err := db.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()
	
	store := From(db)
	
	sm := &entities.SubscriberMetrics{
		UserID:       1,
		Created:      13,
		Deleted:      2,
		Unsubscribed: 6,
		Date:         time.Now(),
	}
	
	err := store.UpdateSubscriberMetrics(sm)
	assert.Nil(t, err)
	
	/*
	We cant test on duplicate key update since the syntax for sqlite is different
	
	sm.Created = 23
	sm.Deleted = 23
	sm.Unsubscribed = 23
	
	err = store.UpdateSubscriberMetrics(sm)
	assert.Nil(t, err)
	*/
}
