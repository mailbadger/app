package storage

import (
	"testing"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/mailbadger/app/entities"
)

func TestScheduledCampaign(t *testing.T) {
	db := openTestDb()
	defer func() {
		err := db.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	var now = time.Now()

	store := From(db)

	c := &entities.Campaign{
		Name:   "foo schedule",
		UserID: 1,
		Status: "draft",
	}

	err := store.CreateCampaign(c)
	assert.Nil(t, err)

	// Test create scheduled campaign
	cs := &entities.CampaignSchedule{
		ID:          ksuid.New(),
		CampaignID:  1,
		ScheduledAt: now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = store.CreateCampaignSchedule(cs)
	assert.Nil(t, err)

	fetchedCampaign, err := store.GetCampaign(c.ID, 1)
	assert.Nil(t, err)
	assert.Equal(t, c.Name, fetchedCampaign.Name)
	assert.Equal(t, entities.StatusScheduled, fetchedCampaign.Status)

	// Test delete scheduled campaign
	err = store.DeleteCampaignSchedule(cs.CampaignID)
	assert.Nil(t, err)

	fetchedCampaign, err = store.GetCampaign(c.ID, 1)
	assert.Nil(t, err)
	assert.Equal(t, c.Name, fetchedCampaign.Name)
	assert.Equal(t, entities.StatusDraft, fetchedCampaign.Status)
}
