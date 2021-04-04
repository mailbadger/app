package storage

import (
	"testing"
	"time"

	"github.com/jinzhu/gorm"
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

	cs, err := store.GetScheduledCampaign(123)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Equal(t, &entities.CampaignSchedules{}, cs)

	// Test create scheduled campaign
	c := &entities.CampaignSchedules{
		ID:          ksuid.New(),
		CampaignID:  1,
		ScheduledAt: now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = store.CreateScheduledCampaign(c)
	assert.Nil(t, err)

	cs, err = store.GetScheduledCampaign(1)
	assert.Nil(t, err)
	assert.Equal(t, c.ID, cs.ID)
	assert.Equal(t, c.CampaignID, cs.CampaignID)
	assert.Equal(t, c.ScheduledAt.UTC(), cs.ScheduledAt.UTC())

	// Test delete scheduled campaign
	err = store.DeleteScheduledCampaign(c.CampaignID)
	assert.Nil(t, err)
}
