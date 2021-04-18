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
	"github.com/mailbadger/app/utils"
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
			continue
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
			continue
		}
		var templateData map[string]string
		err = json.Unmarshal(cs.DefaultTemplateData, &templateData)
		if err != nil {
			continue
		}
		err = template.ValidateData(templateData)
		if err != nil {
			continue
		}

		sesKeys, err := storage.GetSesKeys(c, u.ID)
		if err != nil {
			continue
		}

		segmentIDs, err := utils.StringToIntSlice(cs.SegmentIDs)
		if err != nil {
			continue
		}

		params := &entities.CampaignerTopicParams{
			CampaignID:             cs.CampaignID,
			SegmentIDs:             segmentIDs,
			TemplateData:           templateData,
			Source:                 cs.Source,
			UserID:                 u.ID,
			UserUUID:               u.UUID,
			ConfigurationSetExists: false,
			SesKeys:                *sesKeys,
		}
		paramsByte, err := json.Marshal(params)
		if err != nil {
			continue
		}
		err = queue.Publish(c, entities.CampaignerTopic, paramsByte)
		if err != nil {
			continue
		}
		campaign.Status = entities.StatusSending
		err = storage.UpdateCampaign(c, campaign)
		if err != nil {
			continue
		}
	}

	return nil

}
