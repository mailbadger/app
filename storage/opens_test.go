package storage

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/mailbadger/app/entities"
	"github.com/stretchr/testify/assert"
)

func TestOpens(t *testing.T) {
	db := openTestDb()
	defer func() {
		err := db.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	store := From(db)
	now := time.Now().UTC()

	// Test insert open
	open := []entities.Open{
		{
			ID:         1,
			UserID:     1,
			CampaignID: 1,
			Recipient:  "jhon@doe.com",
			UserAgent:  "android",
			IPAddress:  "1.1.1.1",
			CreatedAt:  now,
		},
		{
			ID:         2,
			UserID:     1,
			CampaignID: 1,
			Recipient:  "jhon@email.com",
			UserAgent:  "windows",
			IPAddress:  "1.1.1.1",
			CreatedAt:  now,
		},
		{
			ID:         3,
			UserID:     1,
			CampaignID: 1,
			Recipient:  "jhon@email.com",
			UserAgent:  "windows",
			IPAddress:  "1.1.1.1",
			CreatedAt:  now,
		},
	}
	// insert open 1
	err := store.CreateOpen(&open[0])
	assert.Nil(t, err)
	// insert open 2
	err = store.CreateOpen(&open[1])
	assert.Nil(t, err)
	// insert open 3
	err = store.CreateOpen(&open[2])
	assert.Nil(t, err)

	// test get campaign opens stats
	opensStats, err := store.GetOpensStats(1)
	assert.Nil(t, err)
	assert.NotNil(t, opensStats)
	assert.Equal(t, int64(3), opensStats.Total)
	assert.Equal(t, int64(2), opensStats.Unique)

}
