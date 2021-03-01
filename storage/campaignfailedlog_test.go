package storage

import (
	"testing"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/mailbadger/app/entities"
)

func TestCampaignFailedLog(t *testing.T) {
	db := openTestDb()
	defer func() {
		err := db.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()
	now := time.Now()

	store := From(db)

	logs := []*entities.CampaignFailedLog{
		{
			UserID:      1,
			CampaignID:  1,
			Description: "reason",
			CreatedAt:   now,
		},
		{
			UserID:      1,
			CampaignID:  2,
			Description: "reason",
			CreatedAt:   now,
		},
		{
			UserID:      1,
			CampaignID:  3,
			Description: "reason",
			CreatedAt:   now,
		},
	}
	campaign1 := &entities.Campaign{
		Model:        entities.Model{ID: 1},
		UserID:       1,
		Name:         "bla",
		TemplateID:   0,
		BaseTemplate: nil,
		Status:       "draft",
		Errors:       nil,
	}
	err := store.CreateCampaign(campaign1)
	assert.Nil(t, err)

	id := ksuid.New()

	// test insert campaign failed log
	for _, k := range logs {
		k.ID = id
		err := store.CreateCampaignFailedLog(k)
		assert.Nil(t, err)
		id = id.Next()
	}

	err = store.LogFailedCampaign(campaign1)
	assert.Nil(t, err)

	c, err := store.GetCampaign(campaign1.ID, campaign1.UserID)
	assert.Nil(t, err)
	assert.Equal(t, "failed", c.Status)
}
