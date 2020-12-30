package storage

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"

	"github.com/mailbadger/app/entities"
)

func createCampaigns(store Storage) {
	for i := 0; i < 100; i++ {
		template := &entities.Template{
			UserID:      1,
			Name:        "bla" + strconv.Itoa(i),
			HTMLPart:    "html_part",
			TextPart:    "text_part",
			SubjectPart: "subject_part",
		}
		err := store.CreateCampaign(&entities.Campaign{
			Name:     "foo " + strconv.Itoa(i),
			Template: template,
			UserID:   1,
			Status:   "draft",
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
		Name: "foo",
		Template: &entities.Template{
			Model: entities.Model{
				ID: 1,
			},
		},
		UserID: 1,
		Status: "draft",
	}

	err := store.CreateCampaign(campaign)
	assert.Nil(t, err)

	//Test get campaign
	campaign, err = store.GetCampaign(campaign.ID, 1)
	assert.Nil(t, err)
	assert.Equal(t, campaign.Name, "foo")
	assert.Equal(t, campaign.Template.ID, 1)

	//Test update campaign
	now := time.Now().UTC()
	campaign.Name = "bar"
	campaign.CompletedAt.SetValid(now)
	err = store.UpdateCampaign(campaign)
	assert.Nil(t, err)
	assert.Equal(t, campaign.Name, "bar")
	assert.True(t, campaign.CompletedAt.Valid)
	assert.Equal(t, campaign.CompletedAt.Time, now)

	//Test get campaigns
	p := NewPaginationCursor("/api/campaigns", 13)
	for i := 0; i < 10; i++ {
		err := store.GetCampaigns(1, p, nil)
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
		err := store.GetCampaigns(1, p, nil)
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

	// Test insert open
	open := []entities.Open{
		{
			ID:         1,
			UserID:     1,
			CampaignID: 1,
			Recipient:  "jhon@doe.com",
			UserAgent:  "android",
			IPAddress:  "1.1.1.1",
			CreatedAt:  now,
		},
		{
			ID:         2,
			UserID:     1,
			CampaignID: 1,
			Recipient:  "jhon@email.com",
			UserAgent:  "windows",
			IPAddress:  "1.1.1.1",
			CreatedAt:  now,
		},
	}

	// test insert open
	for _, i := range open {
		err = store.CreateOpen(&i)
		assert.Nil(t, err)
	}

	//Test get campaign opens backwards
	p = NewPaginationCursor("/api/campaigns/{id}/opens", 13)
	p.SetEndingBefore(1)
	// Test get campaign opens
	err = store.GetCampaignOpens(1, 1, p)
	assert.Nil(t, err)

	campOpens := p.Collection.(*[]entities.Open)
	assert.NotNil(t, *campOpens)
	assert.NotEmpty(t, *campOpens)
	assert.Equal(t, 1, len(*campOpens))
	// campOpens[0] - order desc
	assert.Equal(t, open[1], (*campOpens)[0])

	// insert complaints for test get campaign complaints stats
	complaints := []entities.Complaint{
		{
			ID:         1,
			UserID:     1,
			CampaignID: 1,
			Recipient:  "asd",
			UserAgent:  "dsa",
			Type:       "asd",
			FeedbackID: "dsa",
			CreatedAt:  now,
		},
		{
			ID:         2,
			UserID:     1,
			CampaignID: 1,
			Recipient:  "dsa",
			UserAgent:  "asd",
			Type:       "sda",
			FeedbackID: "w",
			CreatedAt:  now,
		},
	}
	// test insert complaints
	for _, i := range complaints {
		err = store.CreateComplaint(&i)
		assert.Nil(t, err)
	}

	//Test get campaign opens backwards
	p = NewPaginationCursor("/api/campaigns/{id}/complaints", 13)
	p.SetEndingBefore(1)
	// Test get campaign opens
	err = store.GetCampaignComplaints(1, 1, p)
	assert.Nil(t, err)

	campComplaints := p.Collection.(*[]entities.Complaint)
	assert.NotNil(t, *campComplaints)
	assert.NotEmpty(t, *campComplaints)
	assert.Equal(t, 1, len(*campComplaints))
	//  order desc this is why 1 with 0 from slice.
	assert.Equal(t, complaints[1], (*campComplaints)[0])

	// insert bounces to test get campaign bounces
	bounces := []entities.Bounce{
		{
			ID:             1,
			UserID:         1,
			CampaignID:     1,
			Recipient:      "asd",
			Type:           "dsa",
			SubType:        "asd",
			Action:         "dsa",
			Status:         "dsa",
			DiagnosticCode: "dsa",
			FeedbackID:     "dsa",
			CreatedAt:      now,
		},
		{
			ID:             2,
			UserID:         1,
			CampaignID:     1,
			Recipient:      "asd",
			Type:           "dsa",
			SubType:        "asd",
			Action:         "dsa",
			Status:         "dsa",
			DiagnosticCode: "dsa",
			FeedbackID:     "dsa",
			CreatedAt:      now,
		},
	}
	// test insert bounce
	for _, i := range bounces {
		err = store.CreateBounce(&i)
		assert.Nil(t, err)
	}
	//Test get campaign bounces backwards
	p = NewPaginationCursor("/api/campaigns/{id}/bounces", 2)
	p.SetEndingBefore(1)
	// Test get campaign opens
	err = store.GetCampaignBounces(1, 1, p)
	assert.Nil(t, err)

	campBounce := p.Collection.(*[]entities.Bounce)
	assert.NotNil(t, *campBounce)
	assert.NotEmpty(t, *campBounce)
	assert.Equal(t, 1, len(*campBounce))
	// campBounce[0] - order desc
	assert.Equal(t, bounces[1], (*campBounce)[0])

}
