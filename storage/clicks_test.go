package storage

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/mailbadger/app/entities"
)

func TestClicks(t *testing.T) {
	db := openTestDb()
	defer func() {
		err := db.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	store := From(db)
	now := time.Now().UTC()

	// test get empty campaign clicks stats
	clicksStats, err := store.GetClicksStats(1, 1)
	assert.Nil(t, err)
	assert.NotNil(t, clicksStats)
	assert.Equal(t, &entities.ClicksStats{}, clicksStats)

	clicks := []entities.Click{
		{
			ID:         1,
			UserID:     1,
			CampaignID: 1,
			Recipient:  "jhon@doe.com",
			UserAgent:  "android",
			IPAddress:  "1.1.1.1",
			CreatedAt:  now,
			Link:       "s",
		},
		{
			ID:         2,
			UserID:     1,
			CampaignID: 1,
			Recipient:  "jhon@email.com",
			Link:       "a",
			UserAgent:  "windows",
			IPAddress:  "1.1.1.1",
			CreatedAt:  now,
		},
		{
			ID:         3,
			UserID:     1,
			CampaignID: 1,
			Recipient:  "jhon@email.com",
			Link:       "a",
			UserAgent:  "windows",
			IPAddress:  "1.1.1.1",
			CreatedAt:  now,
		},
		{
			ID:         4,
			UserID:     1,
			CampaignID: 2,
			Recipient:  "test@mail.com",
			Link:       "http://example.com?foo=bar",
			UserAgent:  "android",
			IPAddress:  "192.168.0.1",
			CreatedAt:  now,
		},
		{
			ID:         5,
			UserID:     1,
			CampaignID: 2,
			Recipient:  "test@mail.com",
			Link:       "http://example.com?foo=bar",
			UserAgent:  "android",
			IPAddress:  "192.168.0.1",
			CreatedAt:  now,
		},
		{
			ID:         6,
			UserID:     1,
			CampaignID: 2,
			Recipient:  "test@mail.com",
			Link:       "http://example.com?foo=bar",
			UserAgent:  "android",
			IPAddress:  "192.168.0.1",
			CreatedAt:  now,
		},
		{
			ID:         7,
			UserID:     1,
			CampaignID: 2,
			Recipient:  "test@mail.com",
			Link:       "http://example.com?asd=dsa",
			UserAgent:  "android",
			IPAddress:  "192.168.0.1",
			CreatedAt:  now,
		},
		{
			ID:         8,
			UserID:     1,
			CampaignID: 2,
			Recipient:  "test@mail.com",
			Link:       "http://example.com?asd=dsa",
			UserAgent:  "android",
			IPAddress:  "192.168.0.1",
			CreatedAt:  now,
		},
		{
			ID:         9,
			UserID:     1,
			CampaignID: 2,
			Recipient:  "test@mail.com",
			Link:       "http://example.com?test=test",
			UserAgent:  "android",
			IPAddress:  "192.168.0.1",
			CreatedAt:  now,
		},
		{
			ID:         10,
			UserID:     1,
			CampaignID: 2,
			Recipient:  "gl-test@mail.com",
			Link:       "http://example.com?asd=dsa",
			UserAgent:  "windows",
			IPAddress:  "192.168.0.2",
			CreatedAt:  now,
		},
		{
			ID:         11,
			UserID:     1,
			CampaignID: 2,
			Recipient:  "gl-test@mail.com",
			Link:       "http://example.com?foo=bar",
			UserAgent:  "windows",
			IPAddress:  "192.168.0.2",
			CreatedAt:  now,
		},
		{
			ID:         12,
			UserID:     1,
			CampaignID: 2,
			Recipient:  "dj-test@mail.com",
			Link:       "http://example.com?foo=bar",
			UserAgent:  "windows",
			IPAddress:  "192.168.0.2",
			CreatedAt:  now,
		},
		{
			ID:         13,
			UserID:     1,
			CampaignID: 2,
			Recipient:  "dj-test@mail.com",
			Link:       "http://example.com?foo=bar",
			UserAgent:  "windows",
			IPAddress:  "192.168.0.2",
			CreatedAt:  now,
		},
		{
			ID:         14,
			UserID:     1,
			CampaignID: 2,
			Recipient:  "dj-test@mail.com",
			Link:       "http://example.com?test=test",
			UserAgent:  "windows",
			IPAddress:  "192.168.0.2",
			CreatedAt:  now,
		},
	}
	// test insert opens
	for i := range clicks {
		err = store.CreateClick(&clicks[i])
		assert.Nil(t, err)
	}

	// test get campaign clicks stats
	clicksStats, err = store.GetClicksStats(1, 1)
	assert.Nil(t, err)
	assert.NotNil(t, clicksStats)
	exp := &entities.ClicksStats{UniqueClicks: 2, TotalClicks: 3}
	assert.Equal(t, exp, clicksStats)

	// Test get campaign clicks stats
	campaignClicksStats, err := store.GetCampaignClicksStats(2, 1)
	assert.Nil(t, err)
	assert.NotEmpty(t, campaignClicksStats)

	assert.Equal(t, []entities.ClicksStats{
		{
			Link:         "http://example.com?asd=dsa",
			UniqueClicks: 2,
			TotalClicks:  3,
		},
		{
			Link:         "http://example.com?foo=bar",
			UniqueClicks: 3,
			TotalClicks:  6,
		},
		{
			Link:         "http://example.com?test=test",
			UniqueClicks: 2,
			TotalClicks:  2,
		},
	}, campaignClicksStats)

	// Test get empty campaign clicks stats
	campaignClicksStats, err = store.GetCampaignClicksStats(55, 1)
	assert.Nil(t, err)
	assert.Empty(t, campaignClicksStats)
}
