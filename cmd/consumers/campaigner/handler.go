package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/services/campaigns"
	"github.com/mailbadger/app/services/templates"
	awssqs "github.com/mailbadger/app/sqs"
	"github.com/mailbadger/app/storage"
)

type handler struct {
	store             storage.Storage
	campaignsvc       campaigns.Service
	templatesvc       templates.Service
	sqsclient         *sqs.Client
	queueURL          awssqs.CampaignerQueueURL
	sendEmailQueueURL awssqs.SendEmailQueueURL
}

func newHandler(
	store storage.Storage,
	campaignsvc campaigns.Service,
	templatesvc templates.Service,
	sqsclient *sqs.Client,
	queueURL awssqs.CampaignerQueueURL,
	sendEmailQueueURL awssqs.SendEmailQueueURL,
) *handler {
	return &handler{
		store:             store,
		campaignsvc:       campaignsvc,
		templatesvc:       templatesvc,
		sqsclient:         sqsclient,
		queueURL:          queueURL,
		sendEmailQueueURL: sendEmailQueueURL,
	}
}

// HandleMessage is the only requirement needed to fulfill the
// nsq.Handler interface.
func (h *handler) HandleMessage(ctx context.Context, m types.Message) (err error) {
	if m.Body == nil || len(*m.Body) == 0 {
		logrus.Error("Empty message, unable to proceed.")
		return nil
	}

	msg := new(entities.CampaignerTopicParams)
	err = json.Unmarshal([]byte(*m.Body), msg)

	if err != nil {
		logrus.WithError(err).Error("Unable to unmarshal message")
		return err
	}

	logEntry := logrus.WithFields(logrus.Fields{
		"campaign_id": msg.CampaignID,
		"user_id":     msg.UserID,
		"segment_ids": msg.SegmentIDs,
	})

	logEntry.Info("Received a message, processing..")

	campaign, err := h.store.GetCampaign(msg.UserID, msg.CampaignID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logEntry.WithError(err).Warn("campaign does not exist")
			return nil
		}
		logEntry.WithError(err).Error("unable to find campaign")
		return err
	}

	logEntry.WithField("template_id", campaign.TemplateID)

	parsedTemplate, err := h.templatesvc.ParseTemplate(ctx, msg.UserID, campaign.TemplateID)
	if err != nil {
		logEntry.WithError(err).Error("unable to prepare campaign template data")

		err = h.logFailedCampaign(ctx, campaign, "failed to parse template")
		if err != nil {
			logEntry.WithError(err).Errorf("unable to set campaign status to '%s'", entities.StatusFailed)
		}

		return nil
	}

	err = h.processSubscribers(ctx, msg, campaign, parsedTemplate, logEntry, m.ReceiptHandle)
	if err != nil {
		// TODO return wrapped errors and do the logging here instead of inside processSubscribers
		err = h.logFailedCampaign(ctx, campaign, "failed to process subscribers")
		if err != nil {
			logEntry.WithError(err).Errorf("unable to set campaign status to '%s'", entities.StatusFailed)
			return nil
		}

		return nil
	}

	return nil
}

func (h *handler) processSubscribers(
	ctx context.Context,
	msg *entities.CampaignerTopicParams,
	campaign *entities.Campaign,
	parsedTemplate *entities.CampaignTemplateData,
	logEntry *logrus.Entry,
	receiptHandle *string,
) error {
	var (
		timestamp time.Time
		nextID    int64
		limit     int64 = 1000
	)

	id := ksuid.New() // this id will be only used for saving failed send logs

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			subs, err := h.store.GetDistinctSubscribersBySegmentIDs(
				msg.SegmentIDs,
				msg.UserID,
				false, // not in a denylist
				true,  // active
				timestamp,
				nextID,
				limit,
			)
			if err != nil {
				logEntry.WithError(err).Error("unable to fetch subscribers")
				return err
			}
			// extend the timeout for message visibility by 100secs.
			_, err = h.sqsclient.ChangeMessageVisibility(ctx, &sqs.ChangeMessageVisibilityInput{
				QueueUrl:          h.queueURL,
				ReceiptHandle:     receiptHandle,
				VisibilityTimeout: 100,
			})
			if err != nil {
				logrus.WithError(err).Error("unable to extend the message visibility timeout")
			}

			for _, s := range subs {
				id = id.Next()
				params, err := h.campaignsvc.PrepareSubscriberEmailData(
					s,
					*msg,
					campaign.ID,
					parsedTemplate.HTMLPart,
					parsedTemplate.SubjectPart,
					parsedTemplate.TextPart,
				)
				if err != nil {
					logEntry.WithField("subscriber_id", s.ID).WithError(err).Error("unable to prepare subscriber email data")

					err := h.store.CreateSendLog(&entities.SendLog{
						ID:           id,
						UserID:       msg.UserID,
						EventID:      msg.EventID,
						SubscriberID: s.ID,
						CampaignID:   msg.CampaignID,
						Status:       entities.SendLogStatusFailed,
						Description:  fmt.Sprintf("Failed to prepare subscriber email data error: %s", err),
					})
					if err != nil {
						logEntry.WithFields(logrus.Fields{
							"subscriber_id": s.ID,
							"event_id":      msg.EventID,
						}).WithError(err).Error("unable to insert send logs for subscriber.")
					}

					continue
				}

				err = h.campaignsvc.PublishSubscriberEmailParams(ctx, params, h.sendEmailQueueURL)
				if err != nil {
					logEntry.WithField("subscriber_id", s.ID).WithError(err).Error("unable to publish subscriber email params")

					err := h.store.CreateSendLog(&entities.SendLog{
						ID:           id,
						UserID:       msg.UserID,
						EventID:      msg.EventID,
						SubscriberID: s.ID,
						CampaignID:   msg.CampaignID,
						Status:       entities.SendLogStatusFailed,
						Description:  fmt.Sprintf("Failed to publish subscriber email data error: %s", err),
					})
					if err != nil {
						logEntry.WithFields(logrus.Fields{
							"subscriber_id": s.ID,
							"event_id":      msg.EventID,
						}).WithError(err).Error("unable to insert send logs for subscriber.")
					}

					continue
				}
			}

			if len(subs) < 1000 {
				err := h.setStatusSent(ctx, campaign)
				if err != nil {
					logEntry.WithError(err).Errorf("unable to set campaign status to '%s'", entities.StatusSent)
					return err
				}
				return nil
			}

			// set  vars for next batches
			lastSub := subs[len(subs)-1]
			nextID = lastSub.ID
			timestamp = lastSub.CreatedAt
		}
	}
}

// logFailedCampaign updates campaign status to failed & inserts campaign  failed log.
func (h *handler) logFailedCampaign(ctx context.Context, campaign *entities.Campaign, description string) error {
	campaign.Status = entities.StatusFailed
	campaign.CompletedAt.SetValid(time.Now().UTC())
	return h.store.LogFailedCampaign(campaign, description)
}

func (h *handler) setStatusSent(ctx context.Context, campaign *entities.Campaign) error {
	campaign.Status = entities.StatusSent
	campaign.CompletedAt.SetValid(time.Now().UTC())
	return h.store.UpdateCampaign(campaign)
}

func (h *handler) DeleteMessage(ctx context.Context, m types.Message) error {
	_, err := h.sqsclient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      h.queueURL,
		ReceiptHandle: m.ReceiptHandle,
	})
	return err
}
