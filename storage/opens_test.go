package storage

import (
	"testing"
	"time"

	"github.com/mailbadger/app/entities"
	"github.com/stretchr/testify/assert"
)

func TestOpens(t *testing.T) {
	db := openTestDb()

	store := From(db)
	now := time.Now().UTC()

	// test get empty opens stats
	opensStats, err := store.GetOpensStats(1, 1)
	assert.Nil(t, err)
	assert.Equal(t, &entities.OpensStats{}, opensStats)

	// Test insert open
	opens := []entities.Open{
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
			CampaignID: 2,
			Recipient:  "jhon@email.com",
			UserAgent:  "windows",
			IPAddress:  "1.1.1.1",
			CreatedAt:  now,
		},
	}

	// test insert opens
	for i := range opens {
		err = store.CreateOpen(&opens[i])
		assert.Nil(t, err)
	}

	// test get campaign opens stats
	opensStats, err = store.GetOpensStats(1, 1)
	assert.Nil(t, err)
	assert.NotNil(t, opensStats)
	exp := &entities.OpensStats{Unique: 2, Total: 2}
	assert.Equal(t, exp, opensStats)
}
