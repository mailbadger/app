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

	// Test insert click
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
	// insert click 1
	err := store.CreateClick(&clicks[0])
	assert.Nil(t, err)
	// insert click 2
	err = store.CreateClick(&clicks[1])
	assert.Nil(t, err)
	// insert click 3
	err = store.CreateClick(&clicks[2])
	assert.Nil(t, err)

	// test get campaign clicks stats
	clicksStats := &entities.ClicksStats{}
	clicksStats, err = store.GetClicksStats(1)
	assert.Nil(t, err)
	assert.NotNil(t, clicksStats)
	assert.Equal(t, int64(3), clicksStats.Total)
	assert.Equal(t, int64(2), clicksStats.Unique)

}
