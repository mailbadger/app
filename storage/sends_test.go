package storage

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/mailbadger/app/entities"
	"github.com/stretchr/testify/assert"
)

func TestSends(t *testing.T) {
	db := openTestDb()
	defer func() {
		err := db.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	store := From(db)
	now := time.Now().UTC()

	// test get empty sends
	totalSends, err := store.GetTotalSends(1, 1)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), totalSends)

	sends := []entities.Send{
		{
			ID:               1,
			UserID:           1,
			CampaignID:       1,
			MessageID:        "s",
			Source:           "s",
			SendingAccountID: "s",
			Destination:      "s",
			CreatedAt:        now,
		},
		{
			ID:               2,
			UserID:           1,
			CampaignID:       1,
			MessageID:        "a",
			Source:           "a",
			SendingAccountID: "a",
			Destination:      "a",
			CreatedAt:        now,
		},
	}
	// test insert opens
	for i := range sends {
		err = store.CreateSend(&sends[i])
		assert.Nil(t, err)
	}

	// test get total sends stats
	totalSends, err = store.GetTotalSends(1, 1)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), totalSends)

	// test delete all sends for user
	err = store.DeleteAllSendsForUser(1)
	assert.Nil(t, err)

	totalSends, err = store.GetTotalSends(1, 1)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), totalSends)
}
