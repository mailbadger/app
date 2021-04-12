package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/queue"
	"github.com/mailbadger/app/storage"
)

func main() {
	driver := os.Getenv("DATABASE_DRIVER")
	conf := storage.MakeConfigFromEnv(driver)
	s := storage.New(driver, conf)

	now := time.Now()

	err := job(context.Background(), s, now)
	if err != nil {
		panic(err)
	}

}

func job(c context.Context, s storage.Storage, time time.Time) error {
	scheduledCampaigns, err := s.GetScheduledCampaigns(time)
	if err != nil {
		return fmt.Errorf("failed to get scheduled campaigns: %w", err)
	}

	for _, cs := range scheduledCampaigns {
		campaign, err := s.GetCampaign(1, cs.CampaignID)
		if err != nil {
			panic(err)
		}
		if campaign.Status != entities.StatusDraft {
			continue
		}

		// todo get template

		// todo validate template data

		// todo get ses keys

		// todo get segment ids

		params := &entities.CampaignerTopicParams{
			CampaignID:             0,
			SegmentIDs:             nil,
			TemplateData:           nil,
			Source:                 "",
			UserID:                 0,
			UserUUID:               "",
			ConfigurationSetExists: false,
			SesKeys:                entities.SesKeys{},
		}
		paramsByte, err := json.Marshal(params)
		if err != nil {

		}
		err = queue.Publish(c, entities.CampaignerTopic, paramsByte)
		if err != nil {
			panic(err)
		}
	}

	return nil

}
