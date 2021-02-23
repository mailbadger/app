package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime/trace"
	"syscall"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/nsqio/go-nsq"
	"github.com/pkg/profile"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"

	"github.com/mailbadger/app/consumers"
	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/queue"
	"github.com/mailbadger/app/s3"
	"github.com/mailbadger/app/services/campaigns"
	"github.com/mailbadger/app/services/templates"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/utils"
)

// MessageHandler implements the nsq handler interface.
type MessageHandler struct {
	s           storage.Storage
	templatesvc templates.Service
	p           queue.Producer
}

// HandleMessage is the only requirement needed to fulfill the
// nsq.Handler interface.
func (h *MessageHandler) HandleMessage(m *nsq.Message) (err error) {
	ctx, task := trace.NewTask(context.Background(), "handleMessage")
	defer task.End()

	if len(m.Body) == 0 {
		logrus.Error("Empty message, unable to proceed.")
		return nil
	}

	svc := campaigns.New(h.s, h.p)

	msg := new(entities.CampaignerTopicParams)
	trace.WithRegion(ctx, "unmarshalBody", func() {
		err = json.Unmarshal(m.Body, msg)
	})
	if err != nil {
		logrus.WithField("body", string(m.Body)).WithError(err).Error("malformed JSON message")
		return nil
	}

	logEntry := logrus.WithFields(logrus.Fields{
		"campaign_id": msg.CampaignID,
		"user_id":     msg.UserID,
	})

	var campaign *entities.Campaign
	trace.WithRegion(ctx, "getCampaign", func() {
		campaign, err = h.s.GetCampaign(msg.CampaignID, msg.UserID)
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logEntry.WithError(err).Warn("unable to find campaign")
			return nil
		}
		logEntry.WithError(err).Error("unable to find campaign")
		return err
	}
	if campaign.Status != entities.StatusDraft {
		logEntry.WithError(err).Errorf("potentially duplicate message: campaign status is '%s', it should be 'draft'", campaign.Status)
		return nil
	}

	trace.WithRegion(ctx, "setStatusSending", func() {
		campaign.StartedAt.SetValid(time.Now().UTC())
		campaign.Status = entities.StatusSending
		err = h.s.UpdateCampaign(campaign)
	})
	if err != nil {
		logEntry.WithError(err).Errorf("unable to set campaign status to '%s'", entities.StatusSending)
		return err
	}

	var parsedTemplate *entities.CampaignTemplateData
	trace.WithRegion(ctx, "parseTemplate", func() {
		parsedTemplate, err = h.templatesvc.ParseTemplate(context.Background(), campaign.TemplateID, msg.UserID)
	})
	if err != nil {
		logEntry.WithField("template_id", campaign.TemplateID).WithError(err).Error("unable to prepare campaign template data")

		trace.WithRegion(ctx, "setStatusFailed", func() {
			campaign.Status = entities.StatusFailed
			err = h.s.UpdateCampaign(campaign)
		})
		if err != nil {
			logEntry.WithError(err).Errorf("unable to set campaign status to '%s'", entities.StatusFailed)
		}
		return nil
	}

	trace.WithRegion(ctx, "processSubscribers", func() {
		err = processSubscribers(msg, campaign, parsedTemplate, h.s, svc, logEntry)
	})
	if err != nil {
		//TODO return wrapped errors and do the logging here instead of inside processSubscribers
		trace.WithRegion(ctx, "setStatusFailed", func() {
			campaign.Status = entities.StatusFailed
			campaign.CompletedAt.SetValid(time.Now().UTC())
			err = h.s.UpdateCampaign(campaign)
		})
		if err != nil {
			logEntry.WithError(err).Errorf("unable to set campaign status to '%s'", entities.StatusFailed)
			return nil
		}
		return nil
	}

	trace.WithRegion(ctx, "setStatusSent", func() {
		campaign.Status = entities.StatusSent
		campaign.CompletedAt.SetValid(time.Now().UTC())
		err = h.s.UpdateCampaign(campaign)
	})
	if err != nil {
		logEntry.WithError(err).Errorf("unable to set campaign status to '%s'", entities.StatusSent)
		return nil
	}

	return nil
}

func processSubscribers(
	msg *entities.CampaignerTopicParams,
	campaign *entities.Campaign,
	parsedTemplate *entities.CampaignTemplateData,
	store storage.Storage,
	svc campaigns.Service,
	logEntry *logrus.Entry,
) error {
	var (
		timestamp time.Time
		nextID    int64
		limit     int64 = 1000
	)

	id := ksuid.New()

	for {
		subs, err := store.GetDistinctSubscribersBySegmentIDs(
			msg.SegmentIDs,
			msg.UserID,
			false,
			true,
			timestamp,
			nextID,
			limit,
		)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logEntry.WithField("segment_ids", msg.SegmentIDs).WithError(err).Warn("unable to fetch subscribers")
				return err
			}
			logEntry.WithField("segment_ids", msg.SegmentIDs).WithError(err).Error("unable to fetch subscribers")
			return err
		}

		for _, s := range subs {
			params, err := svc.PrepareSubscriberEmailData(s, id, *msg, campaign.ID, parsedTemplate.HTMLPart, parsedTemplate.SubjectPart, parsedTemplate.TextPart)
			if err != nil {
				sendLog := &entities.SendLog{
					ID:           id,
					UserID:       msg.UserID,
					SubscriberID: s.ID,
					CampaignID:   msg.CampaignID,
					Status:       entities.SendLogStatusFailed,
					Description:  fmt.Sprintf("Failed to prepare subscriber email data error: %s", err),
				}
				err := store.CreateSendLog(sendLog)
				if err != nil {
					logEntry.WithFields(logrus.Fields{
						"id":            id.String(),
						"segment_ids":   msg.SegmentIDs,
						"subscriber_id": s.ID,
					}).WithError(err).Error("unable to insert send logs for subscriber.")
				}
				return err
			}

			err = svc.PublishSubscriberEmailParams(params)
			if err != nil {
				sendLog := &entities.SendLog{
					ID:           id,
					UserID:       msg.UserID,
					SubscriberID: s.ID,
					CampaignID:   msg.CampaignID,
					Status:       entities.SendLogStatusFailed,
					Description:  fmt.Sprintf("Failed to publish subscriber email data error: %s", err),
				}
				err := store.CreateSendLog(sendLog)
				if err != nil {
					logEntry.WithFields(logrus.Fields{
						"id":            id.String(),
						"segment_ids":   msg.SegmentIDs,
						"subscriber_id": s.ID,
					}).WithError(err).Error("unable to insert send logs for subscriber.")
				}
				return err
			}

			id = id.Next()
		}

		if len(subs) < 1000 {
			break
		}

		// set  vars for next batches
		lastSub := subs[len(subs)-1]
		nextID = lastSub.ID
		timestamp = lastSub.CreatedAt
	}

	return nil
}

func main() {
	mode := flag.String("profile.mode", "", "enable profiling mode, one of [cpu, mem, mutex, block, trace]")
	flag.Parse()
	switch *mode {
	case "cpu":
		defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	case "mem":
		defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()
	case "mutex":
		defer profile.Start(profile.MutexProfile, profile.ProfilePath(".")).Stop()
	case "block":
		defer profile.Start(profile.BlockProfile, profile.ProfilePath(".")).Stop()
	case "trace":
		defer profile.Start(profile.TraceProfile, profile.ProfilePath(".")).Stop()
	default:
		// do nothing
	}

	lvl, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		lvl = logrus.InfoLevel
	}

	logrus.SetLevel(lvl)
	if utils.IsProductionMode() {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	logrus.SetOutput(os.Stdout)

	driver := os.Getenv("DATABASE_DRIVER")
	conf := storage.MakeConfigFromEnv(driver)
	s := storage.New(driver, conf)

	s3Client, err := s3.NewS3Client(
		os.Getenv("AWS_S3_ACCESS_KEY"),
		os.Getenv("AWS_S3_SECRET_KEY"),
		os.Getenv("AWS_S3_REGION"),
	)
	if err != nil {
		logrus.Fatal(err)
	}

	svc := templates.New(s, s3Client)

	p, err := queue.NewProducer(os.Getenv("NSQD_HOST"), os.Getenv("NSQD_PORT"))
	if err != nil {
		logrus.Fatal(err)
	}

	config := nsq.NewConfig()

	consumer, err := nsq.NewConsumer(entities.CampaignerTopic, entities.CampaignerTopic, config)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create consumer")
	}

	consumer.ChangeMaxInFlight(200)

	consumer.SetLogger(
		&consumers.NoopLogger{},
		nsq.LogLevelError,
	)

	consumer.AddConcurrentHandlers(
		&MessageHandler{s, svc, p},
		20,
	)

	addr := fmt.Sprintf("%s:%s", os.Getenv("NSQLOOKUPD_HOST"), os.Getenv("NSQLOOKUPD_PORT"))
	nsqlds := []string{addr}

	logrus.Infoln("Connecting to NSQlookup...")
	if err := consumer.ConnectToNSQLookupds(nsqlds); err != nil {
		logrus.Fatal(err)
	}

	logrus.Infoln("Connected to NSQlookup")

	shutdown := make(chan os.Signal, 2)
	signal.Notify(shutdown, os.Interrupt)
	signal.Notify(shutdown, syscall.SIGINT)
	signal.Notify(shutdown, syscall.SIGTERM)

	for {
		select {
		case <-consumer.StopChan:
			return // consumer disconnected. Time to quit.
		case <-shutdown:
			// Synchronously drain the queue before falling out of main
			logrus.Infoln("Stopping consumer...")
			consumer.Stop()
		}
	}
}
