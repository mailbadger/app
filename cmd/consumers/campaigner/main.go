package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/jinzhu/gorm"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/mode"
	"github.com/mailbadger/app/s3"
	"github.com/mailbadger/app/services/campaigns"
	"github.com/mailbadger/app/services/templates"
	awssqs "github.com/mailbadger/app/sqs"
	"github.com/mailbadger/app/storage"
)

// MessageHandler implements the nsq handler interface.
type MessageHandler struct {
	store       storage.Storage
	campaignsvc campaigns.Service
	templatesvc templates.Service
	//these are needed to change message visibility
	sqsclient *sqs.Client
	queueURL  *string
}

// HandleMessage is the only requirement needed to fulfill the
// nsq.Handler interface.
func (h *MessageHandler) HandleMessage(ctx context.Context, m types.Message) (err error) {
	if m.Body == nil || len(*m.Body) == 0 {
		logrus.Error("Empty message, unable to proceed.")
		return nil
	}

	msg := new(entities.CampaignerTopicParams)
	err = json.Unmarshal([]byte(*m.Body), msg)

	if err != nil {
		return err
	}

	logEntry := logrus.WithFields(logrus.Fields{
		"campaign_id": msg.CampaignID,
		"user_id":     msg.UserID,
		"segment_ids": msg.SegmentIDs,
	})

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

func (h *MessageHandler) processSubscribers(
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
							"event_id":      msg.EventID.String(),
						}).WithError(err).Error("unable to insert send logs for subscriber.")
					}

					continue
				}

				err = h.campaignsvc.PublishSubscriberEmailParams(ctx, params, h.queueURL)
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
							"event_id":      msg.EventID.String(),
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
			}

			// set  vars for next batches
			lastSub := subs[len(subs)-1]
			nextID = lastSub.ID
			timestamp = lastSub.CreatedAt
		}
	}
}

// logFailedCampaign updates campaign status to failed & inserts campaign  failed log.
func (h *MessageHandler) logFailedCampaign(ctx context.Context, campaign *entities.Campaign, description string) error {
	campaign.Status = entities.StatusFailed
	campaign.CompletedAt.SetValid(time.Now().UTC())
	return h.store.LogFailedCampaign(campaign, description)
}

func (h *MessageHandler) setStatusSent(ctx context.Context, campaign *entities.Campaign) error {
	campaign.Status = entities.StatusSent
	campaign.CompletedAt.SetValid(time.Now().UTC())
	return h.store.UpdateCampaign(campaign)
}

func (h *MessageHandler) DeleteMessage(ctx context.Context, m types.Message) error {
	_, err := h.sqsclient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      h.queueURL,
		ReceiptHandle: m.ReceiptHandle,
	})
	return err
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mode.SetModeFromEnv()

	lvl, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		lvl = logrus.InfoLevel
	}

	logrus.SetLevel(lvl)
	if mode.IsProd() {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	logrus.SetOutput(os.Stdout)

	queueStr := flag.String("q", "", "The name of the queue")
	timeout := flag.Int("t", 360, "How long, in seconds, that the message is hidden from others")
	maxInFlightMsgs := flag.Int("m", 10, "Max number of messages to be received from SQS simultaneously")
	waitTimeout := flag.Int("w", 10, "How long, in seconds, ")
	flag.Parse()

	if *queueStr == "" {
		logrus.Fatal("You must supply the name of a queue (-q QUEUE)")
	}

	if *timeout < 0 {
		*timeout = 0
	}

	if *timeout > 12*60*60 {
		*timeout = 12 * 60 * 60
	}

	if *waitTimeout < 0 {
		*waitTimeout = 0
	}

	if *waitTimeout > 20 {
		*waitTimeout = 20
	}

	if *maxInFlightMsgs < 1 {
		*maxInFlightMsgs = 1
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		logrus.WithError(err).Fatal("AWS configuration error")
	}

	client := sqs.NewFromConfig(cfg)

	gQInput := &sqs.GetQueueUrlInput{
		QueueName: queueStr,
	}
	// Get URL of queue
	urlResult, err := client.GetQueueUrl(ctx, gQInput)
	if err != nil {
		logrus.WithError(err).Fatal("Got an error getting the queue URL")
	}

	queueURL := urlResult.QueueUrl
	consumer := awssqs.NewConsumer(*queueURL, int32(*timeout), int32(*maxInFlightMsgs), int32(*waitTimeout), client)

	driver := os.Getenv("DATABASE_DRIVER")
	conf := storage.MakeConfigFromEnv(driver)
	store := storage.New(driver, conf)

	s3Client, err := s3.NewS3Client(
		os.Getenv("AWS_ACCESS_KEY_ID"),
		os.Getenv("AWS_SECRET_ACCESS_KEY"),
		os.Getenv("AWS_REGION"),
	)
	if err != nil {
		logrus.WithError(err).Fatal("AWS S3 client configuration error")
	}

	templatesvc := templates.New(store, s3Client)
	campaignsvc := campaigns.New(store, client)

	g := new(errgroup.Group)
	handler := &MessageHandler{
		store:       store,
		templatesvc: templatesvc,
		campaignsvc: campaignsvc,
		sqsclient:   client,
		queueURL:    queueURL,
	}
	fn := func(ctx context.Context, m types.Message) func() error {
		return func() error {
			err = handler.HandleMessage(ctx, m)
			if err != nil {
				return err
			}
			return handler.DeleteMessage(ctx, m)
		}
	}

	// Poll for messages
	messages := consumer.PollSQS(ctx)

	for m := range messages {
		g.Go(fn(ctx, m))
	}

	if err := g.Wait(); err != nil {
		logrus.WithError(err).Error("received an error when handling a message")
	}
}
