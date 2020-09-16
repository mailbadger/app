package storage

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/mailbadger/app/entities"
	"github.com/stretchr/testify/assert"
)

func TestClicks(t *testing.T) {
	db := openTestDb()
	defer func() {
		err := db.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	store := From(db)
	now := time.Now().UTC()

	// test get empty campaign clicks stats
	clicksStats, err := store.GetClicksStats(1)
	assert.Nil(t, err)
	assert.NotNil(t, clicksStats)
	assert.Equal(t, &entities.ClicksStats{}, clicksStats)

	clicks := []entities.Click{
		{
			ID:         1,
			UserID:     1,
			CampaignID: 1,
			Recipient:  "jhon@doe.com",
			UserAgent:  "android",
			IPAddress:  "1.1.1.1",
			CreatedAt:  now,
			Link:       "s",
		},
		{
			ID:         2,
			UserID:     1,
			CampaignID: 1,
			Recipient:  "jhon@email.com",
			Link:       "a",
			UserAgent:  "windows",
			IPAddress:  "1.1.1.1",
			CreatedAt:  now,
		},
		{
			ID:         3,
			UserID:     1,
			CampaignID: 1,
			Recipient:  "jhon@email.com",
			Link:       "a",
			UserAgent:  "windows",
			IPAddress:  "1.1.1.1",
			CreatedAt:  now,
		},
	}
	// test insert opens
	for _, i := range clicks {
		err = store.CreateClick(&i)
		assert.Nil(t, err)
	}

	// test get campaign clicks stats
	clicksStats, err = store.GetClicksStats(1)
	assert.Nil(t, err)
	assert.NotNil(t, clicksStats)
	exp := &entities.ClicksStats{Unique: 2, Total: 3}
	assert.Equal(t, exp, clicksStats)

}
