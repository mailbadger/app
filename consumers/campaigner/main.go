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
	"github.com/mailbadger/app/mode"
	"github.com/mailbadger/app/queue"
	"github.com/mailbadger/app/s3"
	"github.com/mailbadger/app/services/campaigns"
	"github.com/mailbadger/app/services/templates"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/storage/redis"
)

// Cache prefix and duration parameters
const (
	cachePrefix   = "campaigner:"
	cacheDuration = 7 * 24 * time.Hour // & days cache duration
)

// MessageHandler implements the nsq handler interface.
type MessageHandler struct {
	s           storage.Storage
	cache       redis.Storage
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
		"segment_ids": msg.SegmentIDs,
	})

	campaign, err := getCampaign(ctx, h.s, msg.UserID, msg.CampaignID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logEntry.WithError(err).Warn("unable to find campaign")
			return nil
		}
		logEntry.WithError(err).Error("unable to find campaign")
		return err
	}

	campaignerKey := redis.GenCacheKey(cachePrefix, msg.EventID.String())

	exist, err := h.cache.Exists(campaignerKey)
	if err != nil {
		logEntry.WithError(err).Error("Unable to check message existence")
		return err
	}

	if exist {
		logEntry.WithError(err).Info("Message already processed")
		return nil
	}

	if err := h.cache.Set(campaignerKey, []byte("sending"), cacheDuration); err != nil {
		logEntry.WithError(err).Error("Unable to write to cache")
		return err
	}

	logEntry.WithField("template_id", campaign.TemplateID)

	parsedTemplate, err := parseTemplate(ctx, h.templatesvc, msg.UserID, campaign.TemplateID)
	if err != nil {
		logEntry.WithError(err).Error("unable to prepare campaign template data")

		err = logFailedCampaign(ctx, h.s, campaign, "failed to parse template")
		if err != nil {
			logEntry.WithError(err).Errorf("unable to set campaign status to '%s'", entities.StatusFailed)
		}

		return nil
	}

	err = processSubscribers(ctx, m, msg, campaign, parsedTemplate, h.s, svc, logEntry)
	if err != nil {
		// TODO return wrapped errors and do the logging here instead of inside processSubscribers
		err = logFailedCampaign(ctx, h.s, campaign, "failed to process subscribers")
		if err != nil {
			logEntry.WithError(err).Errorf("unable to set campaign status to '%s'", entities.StatusFailed)
			return nil
		}

		return nil
	}

	err = setStatusSent(ctx, h.s, campaign)
	if err != nil {
		logEntry.WithError(err).Errorf("unable to set campaign status to '%s'", entities.StatusSent)
		return nil
	}

	return nil
}

func getCampaign(ctx context.Context, store storage.Storage, userID, campaignID int64) (*entities.Campaign, error) {
	defer trace.StartRegion(ctx, "getCampaign").End()

	return store.GetCampaign(campaignID, userID)
}

func parseTemplate(ctx context.Context, templatesvc templates.Service, userID, templateID int64) (*entities.CampaignTemplateData, error) {
	defer trace.StartRegion(ctx, "parseTemplate").End()

	return templatesvc.ParseTemplate(ctx, templateID, userID)
}

func processSubscribers(
	ctx context.Context,
	m *nsq.Message,
	msg *entities.CampaignerTopicParams,
	campaign *entities.Campaign,
	parsedTemplate *entities.CampaignTemplateData,
	store storage.Storage,
	svc campaigns.Service,
	logEntry *logrus.Entry,
) error {
	defer trace.StartRegion(ctx, "processSubscribers").End()

	var (
		timestamp time.Time
		nextID    int64
		limit     int64 = 1000
	)

	id := ksuid.New() // this id will be only used for saving failed send logs

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
				logEntry.WithError(err).Warn("unable to fetch subscribers")
				return err
			}
			logEntry.WithError(err).Error("unable to fetch subscribers")
			return err
		}
		// reset timeout timer for campaigner.
		m.Touch()

		for _, s := range subs {
			id = id.Next()

			params, err := svc.PrepareSubscriberEmailData(s, id, *msg, campaign.ID, parsedTemplate.HTMLPart, parsedTemplate.SubjectPart, parsedTemplate.TextPart)
			if err != nil {
				logEntry.WithField("subscriber_id", s.ID).WithError(err).Error("unable to prepare subscriber email data")

				sendLog := &entities.SendLog{
					ID:           id,
					UserID:       msg.UserID,
					EventID:      msg.EventID,
					SubscriberID: s.ID,
					CampaignID:   msg.CampaignID,
					Status:       entities.SendLogStatusFailed,
					Description:  fmt.Sprintf("Failed to prepare subscriber email data error: %s", err),
				}

				err := store.CreateSendLog(sendLog)
				if err != nil {
					logEntry.WithFields(logrus.Fields{
						"subscriber_id": s.ID,
						"event_id":      msg.EventID.String(),
					}).WithError(err).Error("unable to insert send logs for subscriber.")
				}

				continue
			}

			err = svc.PublishSubscriberEmailParams(params)
			if err != nil {
				logEntry.WithField("subscriber_id", s.ID).WithError(err).Error("unable to publish subscriber email params")

				sendLog := &entities.SendLog{
					ID:           id,
					UserID:       msg.UserID,
					EventID:      msg.EventID,
					SubscriberID: s.ID,
					CampaignID:   msg.CampaignID,
					Status:       entities.SendLogStatusFailed,
					Description:  fmt.Sprintf("Failed to publish subscriber email data error: %s", err),
				}

				err := store.CreateSendLog(sendLog)
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
			break
		}

		// set  vars for next batches
		lastSub := subs[len(subs)-1]
		nextID = lastSub.ID
		timestamp = lastSub.CreatedAt
	}

	return nil
}

// LogFailedMessage overwriting the callback func for max attempts reached to insert into campaign failed logs.
func (h *MessageHandler) LogFailedMessage(m *nsq.Message) {
	if m == nil {
		return
	}
	msg := new(entities.CampaignerTopicParams)
	err := json.Unmarshal(m.Body, &msg)
	if err != nil {
		logrus.WithField("body", string(m.Body)).
			WithError(err).Error("Malformed JSON message.")
		return
	}
	logrus.WithFields(logrus.Fields{
		"user_id":     msg.UserID,
		"campaign_id": msg.CampaignID,
	}).Error("exceeded max attempts for sending the campaign")

	campaign, err := h.s.GetCampaign(msg.CampaignID, msg.UserID)
	if err != nil {
		logrus.WithFields(logrus.Fields{"campaign_id": msg.CampaignID, "user_id": msg.UserID}).
			WithError(err).Error("Failed to get campaign.")
		return
	}
	err = h.s.LogFailedCampaign(campaign, "Exceeded max attempts for preparing the campaign")
	if err != nil {
		logrus.WithField("campaign_id", campaign.ID).
			WithError(err).Error("Failed to store campaign failed log.")
		return
	}
}

// logFailedCampaign updates campaign status to failed & inserts campaign  failed log.
func logFailedCampaign(ctx context.Context, store storage.Storage, campaign *entities.Campaign, description string) error {
	defer trace.StartRegion(ctx, "setStatusFailed").End()

	campaign.Status = entities.StatusFailed
	campaign.CompletedAt.SetValid(time.Now().UTC())
	return store.LogFailedCampaign(campaign, description)
}

func setStatusSent(ctx context.Context, store storage.Storage, campaign *entities.Campaign) error {
	defer trace.StartRegion(ctx, "setStatusSent").End()

	campaign.Status = entities.StatusSent
	campaign.CompletedAt.SetValid(time.Now().UTC())
	return store.UpdateCampaign(campaign)
}

func main() {
	m := flag.String("profile.mode", "", "enable profiling mode, one of [cpu, mem, mutex, block, trace]")
	flag.Parse()
	switch *m {
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

	driver := os.Getenv("DATABASE_DRIVER")
	conf := storage.MakeConfigFromEnv(driver)
	s := storage.New(driver, conf)

	cache, err := redis.NewRedisStore()
	if err != nil {
		logrus.WithError(err).Fatal("Redis: can't establish connection")
	}

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
		&MessageHandler{
			s:           s,
			cache:       cache,
			templatesvc: svc,
			p:           p,
		},
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
