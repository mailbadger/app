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

	//Test create scheduled campaign
	c := &entities.ScheduledCampaign{
		ID:          ksuid.New(),
		CampaignID:  1,
		ScheduledAt: now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err := store.CreateScheduledCampaign(c)
	assert.Nil(t, err)

	// Test delete scheduled campaign
	err = store.DeleteScheduledCampaign(c.ID, c.CampaignID)
	assert.Nil(t, err)
}
