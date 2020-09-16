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

	// Test insert sends
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
	// insert send 1
	err := store.CreateSend(&sends[0])
	assert.Nil(t, err)
	// insert send 2
	err = store.CreateSend(&sends[1])
	assert.Nil(t, err)

	// test get total sends stats
	totalSends, err := store.GetTotalSends(1)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), totalSends)

}
