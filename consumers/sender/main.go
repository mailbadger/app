package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"

	"github.com/mailbadger/app/consumers"
	"github.com/mailbadger/app/emails"
	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/storage/redis"
	"github.com/mailbadger/app/utils"
)

var (
	ErrInvalidSesKeys = errors.New("invali ses keys")
)

const (
	CacheDuration = 300000 * time.Millisecond
	CharSet = "UTF-8"
)

// MessageHandler implements the nsq handler interface.
type MessageHandler struct {
	storage storage.Storage
	cache   redis.Storage
}

// HandleMessage is the only requirement needed to fulfill the
// nsq.Handler interface.
func (h *MessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		logrus.Error("Empty message, unable to proceed.")
		return nil
	}

	msg := new(entities.SenderTopicParams)
	err := json.Unmarshal(m.Body, msg)
	if err != nil {
		logrus.WithField("body", string(m.Body)).
			WithError(err).Error("Malformed JSON message.")
		return nil
	}
	if !h.cache.Exists(msg.UUID) {
		if err := h.cache.Set(msg.UUID, m.Body, CacheDuration); err != nil {
			logrus.WithFields(logrus.Fields{
				"uuid":          msg.UUID,
				"user_id":       msg.UserID,
				"campaign_id":   msg.CampaignID,
				"subscriber_id": msg.SubscriberID,
			}).WithError(err).Error("Unable to write to cache")
			return err
		}
	}

	sendLog := &entities.SendLog{
		UUID:         msg.UUID,
		UserID:       msg.UserID,
		CampaignID:   msg.CampaignID,
		SubscriberID: msg.SubscriberID,
		Status:       entities.StatusDone,
	}

	resp, err := sendEmail(*msg)
	if err != nil {
		sendLog.Status = entities.StatusFailed

		// First check errors for retrying (returning) they don't need to be inserted in send logs
		// also if the error is retryable delete it from cache
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				logrus.WithFields(logrus.Fields{
					"uuid":          msg.UUID,
					"user_id":       msg.UserID,
					"campaign_id":   msg.CampaignID,
					"subscriber_id": msg.SubscriberID,
				}).WithError(aerr).Error("Unable to send bulk templated email. Message rejected.")
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				logrus.WithFields(logrus.Fields{
					"uuid":          msg.UUID,
					"user_id":       msg.UserID,
					"campaign_id":   msg.CampaignID,
					"subscriber_id": msg.SubscriberID,
				}).WithError(aerr).Error("Unable to send bulk templated email. Domain not verified.")
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				logrus.WithFields(logrus.Fields{
					"uuid":          msg.UUID,
					"user_id":       msg.UserID,
					"campaign_id":   msg.CampaignID,
					"subscriber_id": msg.SubscriberID,
				}).WithError(aerr).Error("Unable to send bulk templated email. Configuration set does not exist.")
			default:
				logrus.WithFields(logrus.Fields{
					"uuid":          msg.UUID,
					"user_id":       msg.UserID,
					"campaign_id":   msg.CampaignID,
					"subscriber_id": msg.SubscriberID,
				}).WithError(aerr).Error("Unable to send templated email. Unknown status code.")
			}
		} else {
			logrus.WithFields(logrus.Fields{
				"uuid":          msg.UUID,
				"user_id":       msg.UserID,
				"campaign_id":   msg.CampaignID,
				"subscriber_id": msg.SubscriberID,
			}).WithError(err).Error("Unable to send templated email.")
		}
	}

	if resp != nil {
		sendLog.MessageID = resp.MessageId
		sendLog.Description = resp.GoString()
	}

	err = h.storage.CreateSendLog(sendLog)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"uuid":             msg.UUID,
			"user_id":          msg.UserID,
			"campaign_id":      msg.CampaignID,
			"subscriber_id":    msg.SubscriberID,
			"send_bulk_status": resp.GoString(),
		}).WithError(err).Error("Unable to add log for sent emails result.")
	}

	return nil
}

func sendEmail(msg entities.SenderTopicParams) (*ses.SendEmailOutput, error) {
	if msg.SesKeys.AccessKey == "" || msg.SesKeys.SecretKey == "" || msg.SesKeys.Region == "" {
		return nil, ErrInvalidSesKeys
	}

	client, err := emails.NewSesSender(msg.SesKeys.AccessKey, msg.SesKeys.SecretKey, msg.SesKeys.Region)
	if err != nil {
		return nil, fmt.Errorf("new sender: %w", err)
	}

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(msg.SubscriberEmail)},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(string(msg.HTMLPart)),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(string(msg.TextPart)),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(string(msg.SubjectPart)),
			},
		},
		Source: aws.String(msg.Source),
		Tags: []*ses.MessageTag{
			{
				Name:  aws.String("campaign_id"),
				Value: aws.String(strconv.FormatInt(msg.CampaignID, 10)),
			},
			{
				Name:  aws.String("user_id"),
				Value: aws.String(msg.UserUUID),
			},
		},
	}

	if msg.ConfigurationSetExists {
		input.ConfigurationSetName = aws.String(emails.ConfigurationSetName)
	}

	return client.SendEmail(input)
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

	cache, err := redis.NewRedisStore()
	if err != nil {
		logrus.Fatal(err)
	}
	config := nsq.NewConfig()

	consumer, err := nsq.NewConsumer(entities.SenderTopic, entities.SenderTopic, config)
	if err != nil {
		logrus.Fatal(err)
	}

	consumer.ChangeMaxInFlight(200)

	consumer.SetLogger(
		&consumers.NoopLogger{},
		nsq.LogLevelError,
	)

	consumer.AddConcurrentHandlers(
		&MessageHandler{
			storage: s,
			cache:   cache,
		},
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
