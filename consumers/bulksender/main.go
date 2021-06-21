package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/nsqio/go-nsq"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"

	"github.com/mailbadger/app/consumers"
	"github.com/mailbadger/app/emails"
	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/mode"
	"github.com/mailbadger/app/storage"
)

// MessageHandler implements the nsq handler interface.
type MessageHandler struct {
	s storage.Storage
}

// HandleMessage is the only requirement needed to fulfill the
// nsq.Handler interface.
func (h *MessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		logrus.Error("Empty message, unable to proceed.")
		return nil
	}

	msg := new(entities.BulkSendMessage)

	err := json.Unmarshal(m.Body, msg)
	if err != nil {
		logrus.WithField("body", string(m.Body)).Error("Malformed JSON message.")
		return nil
	}

	if msg.SesKeys == nil {
		logrus.WithField("msg", msg).Error("SES Keys are nil.")
		return nil
	}

	client, err := emails.NewSesSender(msg.SesKeys.AccessKey, msg.SesKeys.SecretKey, msg.SesKeys.Region)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id":     msg.UserID,
			"campaign_id": msg.CampaignID,
		}).WithError(err).Error("Unable to create SES sender")
		return nil
	}

	count, err := h.s.CountLogsByUUID(msg.UUID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id":     msg.UserID,
			"campaign_id": msg.CampaignID,
			"uuid":        msg.UUID,
		}).WithError(err).Error("Unable to count sent logs")
		return nil
	}

	if count > 0 {
		logrus.WithFields(logrus.Fields{
			"user_id":     msg.UserID,
			"campaign_id": msg.CampaignID,
			"uuid":        msg.UUID,
		}).Warn("Bulk already sent.")
		return nil
	}

	res, err := client.SendBulkTemplatedEmail(msg.Input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				logrus.WithFields(logrus.Fields{
					"user_id":     msg.UserID,
					"campaign_id": msg.CampaignID,
					"uuid":        msg.UUID,
				}).WithError(aerr).Error("Unable to send bulk templated email. Message rejected.")
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				logrus.WithFields(logrus.Fields{
					"user_id":     msg.UserID,
					"campaign_id": msg.CampaignID,
					"uuid":        msg.UUID,
				}).WithError(aerr).Error("Unable to send bulk templated email. Domain not verified.")
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				logrus.WithFields(logrus.Fields{
					"user_id":     msg.UserID,
					"campaign_id": msg.CampaignID,
					"uuid":        msg.UUID,
				}).WithError(aerr).Error("Unable to send bulk templated email. Configuration set does not exist.")
			default:
				logrus.WithFields(logrus.Fields{
					"user_id":     msg.UserID,
					"campaign_id": msg.CampaignID,
					"uuid":        msg.UUID,
				}).WithError(aerr).Error("Unable to send bulk templated email. Unknown status code.")
			}
		} else {
			logrus.WithFields(logrus.Fields{
				"user_id":     msg.UserID,
				"campaign_id": msg.CampaignID,
				"uuid":        msg.UUID,
			}).WithError(err).Error("Unable to send bulk templated email.")
		}
		return nil
	}

	id := ksuid.New()

	for _, s := range res.Status {
		err := h.s.CreateSendLog(&entities.SendLog{
			ID:         id,
			UserID:     msg.UserID,
			CampaignID: msg.CampaignID,
			Status:     *s.Status,
		})
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"campaign_id":      msg.CampaignID,
				"user_id":          msg.UserID,
				"send_bulk_status": s.GoString(),
			}).WithError(err).Error("Unable to add log for sent emails result.")
		}

		id = id.Next()
	}

	return nil
}

func main() {
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

	config := nsq.NewConfig()

	consumer, err := nsq.NewConsumer(entities.SendBulkTopic, entities.SendBulkTopic, config)
	if err != nil {
		logrus.Fatal(err)
	}

	consumer.ChangeMaxInFlight(200)

	consumer.SetLogger(
		&consumers.NoopLogger{},
		nsq.LogLevelError,
	)

	consumer.AddConcurrentHandlers(
		&MessageHandler{s},
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
