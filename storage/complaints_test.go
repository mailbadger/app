package storage

import (
	"testing"
	"time"

	"github.com/mailbadger/app/entities"
	"github.com/stretchr/testify/assert"
)

func TestComplaints(t *testing.T) {
	db := openTestDb()

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

	// Test delete all complaints for a user
	err = store.DeleteAllComplaintsForUser(1)
	assert.Nil(t, err)

	totalComplaints, err = store.GetTotalComplaints(1, 1)
	assert.Nil(t, err)
	assert.Empty(t, totalComplaints)
}
