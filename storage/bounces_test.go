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

	// test get empty bounces stats
	totalBounces, err := store.GetTotalBounces(1, 1)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), totalBounces)

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
	// test insert bounces
	for _, i := range bounces {
		err = store.CreateBounce(&i)
		assert.Nil(t, err)
	}

	// test get total bounces stats
	totalBounces, err = store.GetTotalBounces(1, 1)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), totalBounces)

}
