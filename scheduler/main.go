package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/mailbadger/app/emails"
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
	end := time.Since(now)

	logrus.Infof("Scheduler started at %v and took %v minutes to finish", now, end)

}

func job(c context.Context, s storage.Storage, time time.Time) error {
	scheduledCampaigns, err := s.GetScheduledCampaigns(time)
	if err != nil {
		return fmt.Errorf("failed to get scheduled campaigns: %w", err)
	}

	for _, cs := range scheduledCampaigns {

		logEntry := logrus.WithFields(logrus.Fields{
			"campaign_id": cs.CampaignID,
			"user_id":     cs.UserID,
		})

		u, err := s.GetUser(cs.UserID)
		if err != nil {
			logEntry.WithError(err).Error("failed to get user.")
			continue
		}
		campaign, err := s.GetCampaign(u.ID, cs.CampaignID)
		if err != nil {
			logEntry.WithError(err).Error("failed to get campaign.")
			continue
		}
		if campaign.Status != entities.StatusDraft {
			logEntry.WithError(err).Warn("campaign status is not draft.")
			continue
		}

		template, err := storage.GetTemplate(c, campaign.BaseTemplate.ID, u.ID)
		if err != nil {
			logEntry.WithField("template_id", campaign.BaseTemplate.ID).WithError(err).Error("failed to get template.")
			continue
		}
		templateData, err := cs.GetMetadata()
		if err != nil {
			logEntry.WithError(err).Error("failed to unmarshal default template data.")
			continue
		}
		err = template.ValidateData(templateData)
		if err != nil {
			logEntry.WithError(err).Error("failed to validate template data.")
			continue
		}

		sesKeys, err := storage.GetSesKeys(c, u.ID)
		if err != nil {
			logEntry.WithError(err).Error("failed to get ses keys.")
			continue
		}

		segmentIDs, err := cs.GetSegmentIDs()
		if err != nil {
			logEntry.WithError(err).Error("failed to unmarshal segment ids.")
			continue
		}

		lists, err := storage.GetSegmentsByIDs(c, u.ID, segmentIDs)
		if err != nil || len(lists) == 0 {
			logEntry.WithField("segment_ids", segmentIDs).WithError(err).Error("failed to get segments by ids.")
			continue
		}

		sender, err := emails.NewSesSender(sesKeys.AccessKey, sesKeys.SecretKey, sesKeys.Region)
		if err != nil {
			logEntry.WithError(err).Error("failed to create new ses sender.")
			continue
		}

		_, err = sender.DescribeConfigurationSet(&ses.DescribeConfigurationSetInput{
			ConfigurationSetName: aws.String(emails.ConfigurationSetName),
		})

		params := &entities.CampaignerTopicParams{
			CampaignID:             cs.CampaignID,
			SegmentIDs:             segmentIDs,
			TemplateData:           templateData,
			Source:                 cs.Source,
			UserID:                 u.ID,
			UserUUID:               u.UUID,
			ConfigurationSetExists: err == nil,
			SesKeys:                *sesKeys,
		}
		paramsByte, err := json.Marshal(params)
		if err != nil {
			logEntry.WithError(err).Error("failed to marshal params for campaigner.")
			continue
		}
		err = queue.Publish(c, entities.CampaignerTopic, paramsByte)
		if err != nil {
			logEntry.WithError(err).Error("failed to publish campaign to campaigner.")
			continue
		}
		campaign.Status = entities.StatusSending
		err = storage.UpdateCampaign(c, campaign)
		if err != nil {
			logEntry.WithError(err).Error("failed to update status of campaign.")
			continue
		}
	}

	return nil

}
