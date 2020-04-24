package storage

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/news-maily/app/entities"
	"github.com/stretchr/testify/assert"
)

func createCampaigns(store Storage) {
	for i := 0; i < 100; i++ {
		err := store.CreateCampaign(&entities.Campaign{
			Name:         "foo " + strconv.Itoa(i),
			TemplateName: "Template " + strconv.Itoa(i),
			UserID:       1,
			Status:       "draft",
		})
		if err != nil {
			logrus.Fatal(err)
		}
	}
}

func TestCampaign(t *testing.T) {
	db := openTestDb()
	defer func() {
		err := db.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	store := From(db)
	createCampaigns(store)
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
	assert.Equal(t, campaign.Errors["message"], "name: non zero value required")

	//Test get campaigns
	p := NewPaginationCursor("/api/campaigns", 13)
	for i := 0; i < 10; i++ {
		err := store.GetCampaigns(1, p)
		assert.Nil(t, err)
		col := p.Collection.(*[]entities.Campaign)
		assert.NotNil(t, col)
		assert.NotEmpty(t, *col)
		if p.Links.Next != nil {
			assert.Equal(t, int(13), len(*col))
			assert.Equal(t, fmt.Sprintf("/api/campaigns?per_page=%d&starting_after=%d", 13, (*col)[len(*col)-1].GetID()), *p.Links.Next)
			p.SetStartingAfter((*col)[len(*col)-1].GetID())
		} else {
			assert.Equal(t, 10, len(*col))
		}
	}
	assert.Equal(t, int64(101), p.Total)

	//Test get campaigns backwards
	p = NewPaginationCursor("/api/campaigns", 13)
	p.SetEndingBefore(1)
	for i := 0; i < 8; i++ {
		err := store.GetCampaigns(1, p)
		assert.Nil(t, err)
		col := p.Collection.(*[]entities.Campaign)
		assert.NotNil(t, col)
		assert.NotEmpty(t, *col)
		if p.Links.Previous != nil {
			assert.Equal(t, int(13), len(*col))
			assert.Equal(t, fmt.Sprintf("/api/campaigns?ending_before=%d&per_page=%d", (*col)[0].GetID(), 13), *p.Links.Previous)
			p.SetEndingBefore((*col)[0].GetID())
		} else {
			assert.Equal(t, 9, len(*col))
		}
	}

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
