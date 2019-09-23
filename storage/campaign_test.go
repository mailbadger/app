package storage

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/news-maily/app/entities"
	"github.com/news-maily/app/utils/pagination"
	"github.com/stretchr/testify/assert"
)

func TestCampaign(t *testing.T) {
	db := openTestDb()
	defer func() {
		err := db.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	store := From(db)

	//Test create campaign
	campaign := &entities.Campaign{
		Name:         "foo",
		TemplateName: "Template1",
		UserID:       1,
		Status:       "draft",
	}

	err := store.CreateCampaign(campaign)
	assert.Nil(t, err)

	//Test get campaign
	campaign, err = store.GetCampaign(campaign.ID, 1)
	assert.Nil(t, err)
	assert.Equal(t, campaign.Name, "foo")
	assert.Equal(t, campaign.TemplateName, "Template1")

	//Test update campaign
	now := time.Now().UTC()
	campaign.Name = "bar"
	campaign.CompletedAt.SetValid(now)
	err = store.UpdateCampaign(campaign)
	assert.Nil(t, err)
	assert.Equal(t, campaign.Name, "bar")
	assert.True(t, campaign.CompletedAt.Valid)
	assert.Equal(t, campaign.CompletedAt.Time, now)

	//Test campaign validation when name and subject are invalid
	campaign.Name = ""
	campaign.Validate()
	assert.Equal(t, campaign.Errors["name"], entities.ErrCampaignNameEmpty.Error())

	//Test get campaigns
	p := &pagination.Cursor{PerPage: 10}
	store.GetCampaigns(1, p)
	assert.NotEmpty(t, p.Collection)

	//Test get campaigns by template Id
	campaigns, err := store.GetCampaignsByTemplateName("Template1", 1)
	assert.Nil(t, err)
	assert.NotEmpty(t, campaigns)

	//Test get campaign by name
	campaign, err = store.GetCampaignByName("bar", 1)
	assert.Nil(t, err)
	assert.Equal(t, campaign.Name, "bar")

	// Test delete campaign
	err = store.DeleteCampaign(1, 1)
	assert.Nil(t, err)
}
