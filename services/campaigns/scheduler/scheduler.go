package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/mailbadger/app/emails"
	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/logger"
	awssqs "github.com/mailbadger/app/sqs"
	"github.com/mailbadger/app/storage"
	"github.com/sirupsen/logrus"
)

type Scheduler struct {
	s                    storage.Storage
	p                    awssqs.PublisherAPI
	sendCampaignQueueURL awssqs.CampaignerQueueURL
}

func New(
	s storage.Storage,
	p awssqs.PublisherAPI,
	queueURL awssqs.CampaignerQueueURL,
) *Scheduler {
	return &Scheduler{
		s:                    s,
		p:                    p,
		sendCampaignQueueURL: queueURL,
	}
}

func (sched *Scheduler) Start(ctx context.Context, d time.Duration) error {
	logger.From(ctx).Debug("scheduler: starting campaigns scheduler")

	ticker := time.NewTicker(d)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			err := sched.execute(ctx)
			if err != nil {
				logger.From(ctx).WithError(err).Error("scheduler: execute returned error")
			}
		}
	}
}

func (sched *Scheduler) execute(ctx context.Context) error {
	scheduledCampaigns, err := sched.s.GetScheduledCampaigns(time.Now())
	if err != nil {
		return fmt.Errorf("scheduler: failed to get scheduled campaigns: %w", err)
	}

	for _, cs := range scheduledCampaigns {
		logEntry := logrus.WithFields(logrus.Fields{
			"campaign_id": cs.CampaignID,
			"user_id":     cs.UserID,
		})

		u, err := sched.s.GetUser(cs.UserID)
		if err != nil {
			logEntry.WithError(err).Error("sched: failed to get user")
			continue
		}
		campaign, err := sched.s.GetCampaign(cs.CampaignID, u.ID)
		if err != nil {
			logEntry.WithError(err).Error("sched: failed to get campaign")
			continue
		}
		if campaign.Status != entities.StatusScheduled {
			logEntry.WithError(err).Warn("sched: campaign status is not 'scheduled'")
			continue
		}

		template, err := sched.s.GetTemplate(campaign.BaseTemplate.ID, u.ID)
		if err != nil {
			logEntry.WithField("template_id", campaign.BaseTemplate.ID).WithError(err).Error("sched: failed to get template")
			continue
		}
		templateData, err := cs.GetMetadata()
		if err != nil {
			logEntry.WithError(err).Error("sched: failed to unmarshal default template data")
			continue
		}
		err = template.ValidateData(templateData)
		if err != nil {
			logEntry.WithError(err).Error("sched: failed to validate template data")
			continue
		}

		sesKeys, err := sched.s.GetSesKeys(u.ID)
		if err != nil {
			logEntry.WithError(err).Error("sched: failed to get ses keys")
			continue
		}

		segmentIDs, err := cs.GetSegmentIDs()
		if err != nil {
			logEntry.WithError(err).Error("sched: failed to unmarshal segment ids")
			continue
		}

		lists, err := sched.s.GetSegmentsByIDs(u.ID, segmentIDs)
		if err != nil || len(lists) == 0 {
			logEntry.WithField("segment_ids", segmentIDs).WithError(err).Error("sched: failed to get segments by ids")
			continue
		}

		sender, err := emails.NewSesSender(sesKeys.AccessKey, sesKeys.SecretKey, sesKeys.Region)
		if err != nil {
			logEntry.WithError(err).Error("sched: failed to create new ses sender")
			continue
		}

		_, err = sender.DescribeConfigurationSet(&ses.DescribeConfigurationSetInput{
			ConfigurationSetName: aws.String(emails.ConfigurationSetName),
		})

		params := &entities.CampaignerTopicParams{
			EventID:                cs.ID,
			CampaignID:             cs.CampaignID,
			SegmentIDs:             segmentIDs,
			TemplateData:           templateData,
			Source:                 fmt.Sprintf("%s <%s>", cs.FromName, cs.Source),
			UserID:                 u.ID,
			UserUUID:               u.UUID,
			ConfigurationSetExists: err == nil,
			SesKeys:                *sesKeys,
		}
		paramsByte, err := json.Marshal(params)
		if err != nil {
			logEntry.WithError(err).Error("sched: failed to marshal params for campaigner")
			continue
		}
		err = sched.p.SendMessage(ctx, sched.sendCampaignQueueURL, paramsByte)
		if err != nil {
			logEntry.WithError(err).Error("sched: failed to publish campaign to campaigner")
			continue
		}
		campaign.Status = entities.StatusSending
		err = sched.s.UpdateCampaign(campaign)
		if err != nil {
			logEntry.WithError(err).Error("sched: failed to update status of campaign")
			continue
		}
	}

	return nil
}
