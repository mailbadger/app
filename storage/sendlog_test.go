package storage

import (
	"testing"
	"time"

	"github.com/segmentio/ksuid"

	"github.com/stretchr/testify/assert"

	"github.com/mailbadger/app/entities"
)

func TestSendLogs(t *testing.T) {
	db := openTestDb()

	store := From(db)
	now := time.Now().UTC()

	sendLogs := []*entities.SendLog{
		{
			UserID:       1,
			EventID:      ksuid.New(),
			SubscriberID: 1,
			CampaignID:   1,
			Status:       entities.SendLogStatusFailed,
			Description:  "error: some error",
			CreatedAt:    now,
		},
		{
			UserID:       1,
			EventID:      ksuid.New(),
			SubscriberID: 2,
			CampaignID:   1,
			Status:       entities.SendLogStatusFailed,
			Description:  "error: some error",
			CreatedAt:    now,
		},
		{
			UserID:       1,
			EventID:      ksuid.New(),
			SubscriberID: 3,
			CampaignID:   1,
			Status:       entities.SendLogStatusSuccessful,
			Description:  "",
			CreatedAt:    now,
		},
	}

	id := ksuid.New()

	// test insert opens
	for _, sl := range sendLogs {
		sl.ID = id
		err := store.CreateSendLog(sl)
		assert.Nil(t, err)
		id = id.Next()
	}

	n, err := store.CountLogsByStatus(entities.SendLogStatusFailed)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), n)
}
