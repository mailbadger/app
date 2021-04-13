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
		u, err := s.GetUser(cs.UserID)
		if err != nil {
			return err
		}
		campaign, err := s.GetCampaign(u.ID, cs.CampaignID)
		if err != nil {
			panic(err)
		}
		if campaign.Status != entities.StatusDraft {
			continue
		}

		template, err := storage.GetTemplate(c, campaign.BaseTemplate.ID, u.ID)
		if err != nil {
			return err
		}

		// fixme: default template data missing
		err = template.ValidateData(nil)
		if err != nil {
			return err
		}

		sesKeys, err := storage.GetSesKeys(c, u.ID)
		if err != nil {
			return nil
		}

		params := &entities.CampaignerTopicParams{
			CampaignID:             cs.CampaignID,
			SegmentIDs:             nil,
			TemplateData:           nil,
			Source:                 "job_scheduler",
			UserID:                 u.ID,
			UserUUID:               u.UUID,
			ConfigurationSetExists: false,
			SesKeys:                *sesKeys,
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
