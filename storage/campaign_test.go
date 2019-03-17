package storage

import (
	"testing"

	"github.com/news-maily/api/entities"
	"github.com/news-maily/api/utils/pagination"
	"github.com/stretchr/testify/assert"
)

func TestCampaign(t *testing.T) {
	db := openTestDb()
	defer db.Close()

	store := From(db)

	//Test create campaign
	campaign := &entities.Campaign{
		Name:         "foo",
		TemplateName: "Template1",
		UserId:       1,
		Status:       "draft",
	}

	err := store.CreateCampaign(campaign)
	assert.Nil(t, err)

	//Test get campaign
	campaign, err = store.GetCampaign(campaign.Id, 1)
	assert.Nil(t, err)
	assert.Equal(t, campaign.Name, "foo")
	assert.Equal(t, campaign.TemplateName, "Template1")

	//Test update campaign
	campaign.Name = "bar"
	err = store.UpdateCampaign(campaign)
	assert.Nil(t, err)
	assert.Equal(t, campaign.Name, "bar")

	//Test campaign validation when name and subject are invalid
	campaign.Name = ""
	campaign.Validate()
	assert.Equal(t, campaign.Errors["name"], entities.ErrCampaignNameEmpty.Error())

	//Test get campaigns
	p := &pagination.Pagination{PerPage: 10}
	store.GetCampaigns(1, p)
	assert.NotEmpty(t, p.Collection)
	assert.Equal(t, len(p.Collection), int(p.Total))

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
