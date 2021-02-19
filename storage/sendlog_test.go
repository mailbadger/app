package storage

import (
	"testing"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"

	"github.com/mailbadger/app/entities"
)

func TestSendLogs(t *testing.T) {
	db := openTestDb()
	defer func() {
		err := db.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	store := From(db)
	now := time.Now().UTC()

	sendLogs := []*entities.SendLog{
		{
			UserID:       1,
			SubscriberID: 1,
			CampaignID:   1,
			Status:       entities.SendLogStatusFailed,
			Description:  "error: some error",
			CreatedAt:    now,
		},
		{
			UserID:       1,
			SubscriberID: 2,
			CampaignID:   1,
			Status:       entities.SendLogStatusFailed,
			Description:  "error: some error",
			CreatedAt:    now,
		},
		{
			UserID:       1,
			SubscriberID: 3,
			CampaignID:   1,
			Status:       entities.SendLogStatusSuccessful,
			Description:  "",
			CreatedAt:    now,
		},
	}

	uid := ksuid.New()

	// test insert opens
	for _, sl := range sendLogs {
		sl.UID = uid.Next().String()
		err := store.CreateSendLog(sl)
		assert.Nil(t, err)
		uid = uid.Next()
	}

	n, err := store.CountLogsByStatus(entities.SendLogStatusFailed)
	assert.Nil(t, err)
	assert.Equal(t, 2, n)

}
