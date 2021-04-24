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
	"github.com/sirupsen/logrus"
)

func main() {
	driver := os.Getenv("DATABASE_DRIVER")
	conf := storage.MakeConfigFromEnv(driver)
	s := storage.New(driver, conf)

	now := time.Now()
	err := job(context.Background(), s, now)
	if err != nil {
		logrus.WithField("time", now).WithError(err).Error("failed to start campaign scheduler job")
	}
	end := time.Now()

	logrus.Infof("Scheduler started at %s, ended at: %s and took %f minutes to finish", now.String(), end.String(), end.Sub(now).Minutes())

}

func job(c context.Context, s storage.Storage, time time.Time) error {
	scheduledCampaigns, err := s.GetScheduledCampaigns(time)
	if err != nil {
		return fmt.Errorf("failed to get scheduled campaigns: %w", err)
	}

	for _, cs := range scheduledCampaigns {
		u, err := s.GetUser(cs.UserID)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"campaign_id": cs.CampaignID,
				"user_id":     cs.UserID,
			}).WithError(err).Error("failed to get user.")
			continue
		}
		campaign, err := s.GetCampaign(u.ID, cs.CampaignID)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"campaign_id": cs.CampaignID,
				"user_id":     cs.UserID,
			}).WithError(err).Error("failed to get campaign.")
			continue
		}
		if campaign.Status != entities.StatusDraft {
			logrus.WithFields(logrus.Fields{
				"campaign_id": cs.CampaignID,
				"user_id":     cs.UserID,
			}).WithError(err).Warn("campaign status is not draft.")
			continue
		}

		template, err := storage.GetTemplate(c, campaign.BaseTemplate.ID, u.ID)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"campaign_id": cs.CampaignID,
				"user_id":     cs.UserID,
				"template_id": campaign.BaseTemplate.ID,
			}).WithError(err).Error("failed to get template.")
			continue
		}
		var templateData map[string]string
		err = json.Unmarshal(cs.DefaultTemplateData, &templateData)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"campaign_id": cs.CampaignID,
				"user_id":     cs.UserID,
			}).WithError(err).Error("failed to unmarshal default template data.")
			continue
		}
		err = template.ValidateData(templateData)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"campaign_id": cs.CampaignID,
				"user_id":     cs.UserID,
			}).WithError(err).Error("failed to validate template data.")
			continue
		}

		sesKeys, err := storage.GetSesKeys(c, u.ID)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"campaign_id": cs.CampaignID,
				"user_id":     cs.UserID,
			}).WithError(err).Error("failed to get ses keys.")
			continue
		}

		var segmentIDs []int64
		err = json.Unmarshal(cs.SegmentIDs, &segmentIDs)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"campaign_id": cs.CampaignID,
				"user_id":     cs.UserID,
			}).WithError(err).Error("failed to unmarshal segment ids.")
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
			logrus.WithFields(logrus.Fields{
				"campaign_id": cs.CampaignID,
				"user_id":     cs.UserID,
			}).WithError(err).Error("failed to marshal params for campaigner.")
			continue
		}
		err = queue.Publish(c, entities.CampaignerTopic, paramsByte)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"campaign_id": cs.CampaignID,
				"user_id":     cs.UserID,
			}).WithError(err).Error("failed to publish campaign to campaigner.")
			continue
		}
		campaign.Status = entities.StatusScheduled
		err = storage.UpdateCampaign(c, campaign)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"campaign_id": cs.CampaignID,
				"user_id":     cs.UserID,
			}).WithError(err).Error("failed to update status of campaign.")
			continue
		}
	}

	return nil

}
