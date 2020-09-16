package storage

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/mailbadger/app/entities"
	"github.com/stretchr/testify/assert"
)

func TestBounces(t *testing.T) {
	db := openTestDb()
	defer func() {
		err := db.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	store := From(db)
	now := time.Now().UTC()

	// Test insert bounces
	bounces := []entities.Bounce{
		{
			ID:             1,
			UserID:         1,
			CampaignID:     1,
			Recipient:      "jhon@doe.com",
			Type:           "bla",
			SubType:        "bla",
			Action:         "act",
			Status:         "st",
			DiagnosticCode: "asd",
			FeedbackID:     "bla",
			CreatedAt:      now,
		},
		{
			ID:             2,
			UserID:         1,
			CampaignID:     1,
			Recipient:      "jhon@email.com",
			Type:           "bla",
			SubType:        "bla",
			Action:         "act",
			Status:         "st",
			DiagnosticCode: "asd",
			FeedbackID:     "s",
			CreatedAt:      now,
		},
	}
	// insert bounce 1
	err := store.CreateBounce(&bounces[0])
	assert.Nil(t, err)
	// insert bounce 2
	err = store.CreateBounce(&bounces[1])
	assert.Nil(t, err)

	// test get total bounces stats
	totalBounces, err := store.GetTotalBounces(1)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), totalBounces)

}
