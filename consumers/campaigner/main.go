package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/google/uuid"
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"

	"github.com/mailbadger/app/consumers"
	"github.com/mailbadger/app/emails"
	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/queue"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/utils"
)

// MessageHandler implements the nsq handler interface.
type MessageHandler struct {
	s storage.Storage
	p queue.Producer
}

// HandleMessage is the only requirement needed to fulfill the
// nsq.Handler interface.
func (h *MessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		logrus.Error("Empty message, unable to proceed.")
		return nil
	}

	msg := new(entities.CampaignerMessageBody)

	err := json.Unmarshal(m.Body, msg)
	if err != nil {
		logrus.WithField("body", string(m.Body)).WithError(err).Error("Malformed JSON message.")
		return nil
	}

	c, err := h.s.GetCampaign(msg.CampaignID, msg.UserID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"campaign_id": msg.CampaignID,
			"user_id":     msg.UserID,
		}).WithError(err).Error("Unable to fetch campaign")
		return nil
	}

	// fetching subs that are active and that have not been blacklisted
	var (
		active      = true
		blacklisted = false
		timestamp   time.Time
		nextID      int64
		limit       int64 = 1000
	)
	for {
		subs, err := h.s.GetDistinctSubscribersBySegmentIDs(
			msg.SegmentIDs,
			msg.UserID,
			blacklisted,
			active,
			timestamp,
			nextID,
			limit,
		)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"user_id":     msg.UserID,
				"segment_ids": msg.SegmentIDs,
			}).WithError(err).Error("Unable to fetch subscribers.")
			break
		}

		if len(subs) == 0 {
			break
		}

		// SES allows to send 50 emails in a bulk sending operation
		chunkSize := 50
		for i := 0; i < len(subs); i += chunkSize {
			end := i + chunkSize
			if end > len(subs) {
				end = len(subs)
			}

			var dest []*ses.BulkEmailDestination
			for _, s := range subs[i:end] {
				m, err := s.GetMetadata()
				if err != nil {
					logrus.WithError(err).
						WithField("subscriber", s).
						Error("Unable to get subscriber metadata.")

					continue
				}

				if s.Name != "" {
					m["name"] = s.Name
				}

				url, err := s.GetUnsubscribeURL(msg.UserUUID)
				if err != nil {
					logrus.WithError(err).
						WithField("subscriber", s).
						Error("Unable to get unsubscribe url.")
				} else {
					m["unsubscribe_url"] = url
				}

				jsonMeta, err := json.Marshal(m)
				if err != nil {
					logrus.WithError(err).
						WithField("subscriber", s).
						Error("Unable to marshal metadata to json.")

					continue
				}
				d := &ses.BulkEmailDestination{
					Destination: &ses.Destination{
						ToAddresses: []*string{aws.String(s.Email)},
					},
					ReplacementTemplateData: aws.String(string(jsonMeta)),
				}

				dest = append(dest, d)
			}

			defaultData, err := json.Marshal(msg.TemplateData)
			if err != nil {
				logrus.WithError(err).Error("Unable to marshal template data as JSON.")
				continue
			}

			// prepare message for publishing to the queue
			input := &ses.SendBulkTemplatedEmailInput{
				Source:              aws.String(msg.Source),
				Template:            aws.String(c.Template.Name),
				Destinations:        dest,
				DefaultTemplateData: aws.String(string(defaultData)),
				DefaultTags: []*ses.MessageTag{
					&ses.MessageTag{
						Name:  aws.String("campaign_id"),
						Value: aws.String(strconv.FormatInt(msg.CampaignID, 10)),
					},
					&ses.MessageTag{
						Name:  aws.String("user_id"),
						Value: aws.String(msg.UserUUID),
					},
				},
			}

			if msg.ConfigurationSetExists {
				input.ConfigurationSetName = aws.String(emails.ConfigurationSetName)
			}

			uuid := uuid.New()

			bulkMsg := entities.BulkSendMessage{
				UUID:       uuid.String(),
				Input:      input,
				CampaignID: msg.CampaignID,
				UserID:     msg.UserID,
				SesKeys:    &msg.SesKeys,
			}
			msg, err := json.Marshal(bulkMsg)

			if err != nil {
				logrus.WithError(err).Error("Unable to marshal bulk message input.")
				continue
			}

			// publish the message to the queue
			err = h.p.Publish(entities.SendBulkTopic, msg)
			if err != nil {
				logrus.WithError(err).Error("Unable to publish message to send bulk topic.")
			}
		}

		lastSub := subs[len(subs)-1]
		nextID = lastSub.ID
		timestamp = lastSub.CreatedAt
	}

	c.UserID = msg.UserID
	c.Status = entities.StatusSent
	c.CompletedAt.SetValid(time.Now().UTC())

	err = h.s.UpdateCampaign(c)
	if err != nil {
		logrus.WithField("campaign", c).WithError(err).Error("Unable to update campaign.")
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

	p, err := queue.NewProducer(os.Getenv("NSQD_HOST"), os.Getenv("NSQD_PORT"))
	if err != nil {
		logrus.Fatal(err)
	}

	config := nsq.NewConfig()

	consumer, err := nsq.NewConsumer(entities.CampaignerTopic, entities.CampaignerTopic, config)
	if err != nil {
		logrus.Fatal(err)
	}

	consumer.ChangeMaxInFlight(200)

	consumer.SetLogger(
		&consumers.NoopLogger{},
		nsq.LogLevelError,
	)

	consumer.AddConcurrentHandlers(
		&MessageHandler{s, p},
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
