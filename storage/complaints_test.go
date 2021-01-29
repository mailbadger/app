package storage

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/mailbadger/app/entities"
	"github.com/stretchr/testify/assert"
)

func TestComplaints(t *testing.T) {
	db := openTestDb()
	defer func() {
		err := db.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	store := From(db)
	now := time.Now().UTC()

	// Test get empty complaints
	totalComplaints, err := store.GetTotalComplaints(1, 1)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), totalComplaints)

	complaints := []entities.Complaint{
		{
			ID:         1,
			UserID:     1,
			CampaignID: 1,
			Recipient:  "jhon@doe.com",
			UserAgent:  "android",
			Type:       "bla",
			FeedbackID: "bla",
			CreatedAt:  now,
		},
		{
			ID:         2,
			UserID:     1,
			CampaignID: 1,
			Recipient:  "jhon@email.com",
			UserAgent:  "windows",
			Type:       "bla",
			FeedbackID: "bla",
			CreatedAt:  now,
		},
	}
	// test insert opens
	for i := range complaints {
		err = store.CreateComplaint(&complaints[i])
		assert.Nil(t, err)
	}

	// test get total complaints
	totalComplaints, err = store.GetTotalComplaints(1, 1)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), totalComplaints)

}
