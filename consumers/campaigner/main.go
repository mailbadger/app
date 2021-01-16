package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cbroglie/mustache"
	"github.com/google/uuid"
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"

	"github.com/mailbadger/app/consumers"
	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/queue"
	"github.com/mailbadger/app/s3"
	"github.com/mailbadger/app/services/templates"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/utils"
)

// MessageHandler implements the nsq handler interface.
type MessageHandler struct {
	s   storage.Storage
	svc templates.Service
	p   queue.Producer
}

// HandleMessage is the only requirement needed to fulfill the
// nsq.Handler interface.
func (h *MessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		logrus.Error("Empty message, unable to proceed.")
		return nil
	}

	msg := new(entities.SendCampaignParams)

	err := json.Unmarshal(m.Body, msg)
	if err != nil {
		logrus.WithField("body", string(m.Body)).WithError(err).Error("Malformed JSON message.")
		return nil
	}

	campaign, err := h.s.GetCampaign(msg.CampaignID, msg.UserID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"campaign_id": msg.CampaignID,
			"user_id":     msg.UserID},
		).WithError(err).Error("unable to find campaign")
		return nil
	}
	if campaign.Status != entities.StatusDraft {
		logrus.WithError(err).
			WithField("status", campaign.Status).
			Warn("this is duplicate message in campaigner consumer")
		return nil
	}

	campaign.Status = entities.StatusSending
	err = h.s.UpdateCampaign(campaign)
	if err != nil {
		logrus.WithError(err).
			WithField("campaign", campaign).
			Error("unable to update campaign")
		return err
	}

	var (
		timestamp time.Time
		nextID    int64
		limit     int64 = 1000
	)
	for {
		subs, err := h.s.GetDistinctSubscribersBySegmentIDs(
			msg.SegmentIDs,
			msg.UserID,
			false,
			true,
			timestamp,
			nextID,
			limit,
		)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"user_id":     msg.UserID,
				"segment_ids": msg.SegmentIDs,
			}).WithError(err).Error("Unable to fetch subscribers.")
			return nil
		}

		if len(subs) == 0 {
			break
		}

		template, err := h.svc.GetTemplate(context.Background(), msg.TemplateID, msg.UserID)
		if err != nil {
			logrus.WithError(err).
				WithField("template_id", msg.TemplateID).
				Error("unable to get template")
			return nil
		}
		html, err := mustache.ParseString(template.HTMLPart)
		if err != nil {
			logrus.WithError(err).
				WithField("html_part", template.HTMLPart).
				Error("unable to parse html_part")
			return nil
		}
		txt, err := mustache.ParseString(template.SubjectPart)
		if err != nil {
			logrus.WithError(err).
				WithField("subject_part", template.SubjectPart).
				Error("unable to parse subject_part")
			return nil
		}
		sub, err := mustache.ParseString(template.TextPart)
		if err != nil {
			logrus.WithError(err).
				WithField("text_part", template.TextPart).
				Error("unable to parse text_part")
			return nil
		}

		for _, s := range subs {
			m, err := s.GetMetadata()
			if err != nil {
				logrus.WithError(err).
					WithField("subscriber", s).
					Error("Unable to get subscriber metadata.")
				continue
			}
			//  todo merge def template data with sub metadata

			// fill template buffer with rendered template.
			var (
				htmlBuf bytes.Buffer
				subBuf  bytes.Buffer
				textBuf bytes.Buffer
			)

			err = html.FRender(&htmlBuf, m)
			if err != nil {
				logrus.WithError(err).Error("unable to render html_part")
			}
			err = txt.FRender(&subBuf, m)
			if err != nil {
				logrus.WithError(err).
					WithField("subject_part", template.SubjectPart).
					Error("unable to render text_part")
			}
			err = sub.FRender(&textBuf, m)
			if err != nil {
				logrus.WithError(err).
					WithField("subject_part", template.SubjectPart).
					Error("unable to render subject_part")
			}

			sender := entities.SendEmailTopicParams{
				UUID:         uuid.New().String(),
				SubscriberID: s.ID,
				CampaignID:   campaign.ID,
				SesKeys:      msg.SesKeys,
				HTMLPart:     htmlBuf.Bytes(),
				SubjectPart:  subBuf.Bytes(),
				TextPart:     textBuf.Bytes(),
				UserUUID:     msg.UserUUID,
				UserID:       msg.UserID,
			}

			senderBytes, err := json.Marshal(sender)
			if err != nil {
				logrus.WithError(err).Error("Unable to marshal bulk message input.")
				continue
			}

			// publish the message to the queue
			err = h.p.Publish(entities.SenderTopic, senderBytes)
			if err != nil {
				logrus.WithError(err).Error("Unable to publish message to send bulk topic.")
				continue
			}
		}

		// set  vars for next batches
		lastSub := subs[len(subs)-1]
		nextID = lastSub.ID
		timestamp = lastSub.CreatedAt

		// exit the loop if next batch is smaller then 1k -> it's the last batch of 1k.
		if len(subs) < 1000 {
			continue
		}
	}

	campaign.Status = entities.StatusSent
	campaign.CompletedAt.SetValid(time.Now().UTC())
	err = h.s.UpdateCampaign(campaign)
	if err != nil {
		logrus.WithError(err).
			WithField("campaign", campaign).
			Error("unable to update campaign")
		return err
	}

	return nil
}

func main() {
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

	consumer, err := nsq.NewConsumer(entities.CampaignsTopic, entities.CampaignsTopic, config)
	if err != nil {
		logrus.Fatal(err)
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
