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
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/nsqio/go-nsq"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"

	"github.com/mailbadger/app/consumers"
	"github.com/mailbadger/app/emails"
	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/mode"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/storage/redis"
)

// Sender errors
var (
	ErrInvalidSesKeys = errors.New("invalid ses keys")
)

// Cache prefix and duration parameters
const (
	cachePrefix   = "sender:"
	cacheDuration = 7 * 24 * time.Hour // & days cache duration
)

// CharSet is used for the SES message body charset
const CharSet = "UTF-8"

// MessageHandler implements the nsq handler interface.
type MessageHandler struct {
	storage storage.Storage
	cache   redis.Storage
}

// LogFailedMessage is for overriding the nsq.FailedMessageLogger
// interface which handles the last failing retry
func (h *MessageHandler) LogFailedMessage(m *nsq.Message) {
	if m == nil || len(m.Body) == 0 {
		logrus.Error("Empty message, unable to proceed with failed message.")
		return
	}

	msg := new(entities.SenderTopicParams)
	err := json.Unmarshal(m.Body, msg)
	if err != nil {
		logrus.WithField("body", string(m.Body)).
			WithError(err).Error("Malformed JSON message.")
		return
	}

	logEntry := logrus.WithFields(logrus.Fields{
		"event_id":      msg.EventID,
		"user_id":       msg.UserID,
		"campaign_id":   msg.CampaignID,
		"subscriber_id": msg.SubscriberID,
	})

	logEntry.Error("Exceeded max attempts for sending the e-mail.")

	err = h.storage.CreateSendLog(&entities.SendLog{
		ID:           ksuid.New(),
		EventID:      msg.EventID,
		UserID:       msg.UserID,
		CampaignID:   msg.CampaignID,
		SubscriberID: msg.SubscriberID,
		Status:       entities.SendLogStatusFailed,
		Description:  "Exceeded max attempts for sending the e-mail.",
	})
	if err != nil {
		logEntry.WithError(err).Error("Unable to add log for sent emails result.")
	}
}

// HandleMessage is the only requirement needed to fulfill the
// nsq.Handler interface.
func (h *MessageHandler) HandleMessage(m *nsq.Message) (err error) {
	if len(m.Body) == 0 {
		logrus.Error("Empty message, unable to proceed.")
		return nil
	}

	msg := new(entities.SenderTopicParams)
	err = json.Unmarshal(m.Body, msg)
	if err != nil {
		logrus.WithField("body", string(m.Body)).
			WithError(err).Error("Malformed JSON message.")
		return nil
	}

	logEntry := logrus.WithFields(logrus.Fields{
		"event_id":      msg.EventID,
		"user_id":       msg.UserID,
		"campaign_id":   msg.CampaignID,
		"subscriber_id": msg.SubscriberID,
	})

	cacheKey := redis.GenCacheKey(cachePrefix, fmt.Sprintf("%s_%d", msg.EventID.String(), msg.SubscriberID))

	// check if the message is processing (if the uuid exists in redis that means it is in progress)
	exist, err := h.cache.Exists(cacheKey)
	if err != nil {
		logEntry.WithError(err).Error("Unable to check message existence")
		return err
	}

	if exist {
		logEntry.WithError(err).Info("Message already processed")
		return nil
	}

	if err := h.cache.Set(cacheKey, []byte("sending"), cacheDuration); err != nil {
		logEntry.WithError(err).Error("Unable to write to cache")
		return err
	}

	sendLog := &entities.SendLog{
		ID:           ksuid.New(),
		EventID:      msg.EventID,
		UserID:       msg.UserID,
		CampaignID:   msg.CampaignID,
		SubscriberID: msg.SubscriberID,
		Status:       entities.SendLogStatusSuccessful,
		Description:  entities.SendLogDescriptionOnSuccessful,
	}

	defer func() {
		if err == nil {
			err = h.storage.CreateSendLog(sendLog)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"event_id":      msg.EventID.String(),
					"user_id":       msg.UserID,
					"campaign_id":   msg.CampaignID,
					"subscriber_id": msg.SubscriberID,
					"message_id":    sendLog.MessageID,
				}).WithError(err).Error("Unable to add log for sent emails result.")
			}
		}
	}()

	client, err := newSesClient(msg.SesKeys)
	if err != nil {
		logEntry.WithError(err).Error("Unable to create ses sender")

		sendLog.Status = entities.StatusFailed
		sendLog.Description = entities.SendLogDescriptionOnSesClientError

		return nil
	}

	resp, err := sendEmail(client, *msg)
	if err != nil {
		sendLog.Status = entities.StatusFailed
		sendLog.Description = entities.SendLogDescriptionOnSendEmailError

		// First check errors for retrying (returning) they don't need to be inserted in send logs
		// also if the error is retryable delete it from cache
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				sendLog.Description = "Unable to send email, message rejected."
				logEntry.WithError(aerr).Error("Unable to send templated email. Message rejected.")
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				sendLog.Description = "Unable to send email, domain not verified."
				logEntry.WithError(aerr).Error("Unable to send templated email. Domain not verified.")
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				sendLog.Description = "Unable to send email, configuration set does not exist."
				logEntry.WithError(aerr).Error("Unable to send templated email. Configuration set does not exist.")
			case sns.ErrCodeThrottledException:
				logEntry.WithError(aerr).Error("Unable to send templated email. The rate at which requests have been submitted for this action exceeds the limit for your account. Slow down!")
				rerr := h.cache.Delete(cacheKey)
				if rerr != nil {
					logEntry.WithError(rerr).Error("Unable to delete cached id")
				}
				return err
			case sns.ErrCodeInternalErrorException:
				logEntry.WithError(aerr).Error("Unable to send templated email. The request processing has failed because of an unknown error, exception, or failure.")
				rerr := h.cache.Delete(cacheKey)
				if rerr != nil {
					logEntry.WithError(rerr).Error("Unable to delete cached id")
				}
				return err
			default:
				logEntry.WithError(aerr).Error("Unable to send templated email. Unknown status.")
				rerr := h.cache.Delete(cacheKey)
				if rerr != nil {
					logEntry.WithError(rerr).Error("Unable to delete cached id")
				}
				return err
			}
		} else {
			logEntry.WithError(err).Error("Unable to send templated email.")
			rerr := h.cache.Delete(cacheKey)
			if rerr != nil {
				logEntry.WithError(rerr).Error("Unable to delete cached id")
			}
			return err
		}
	}

	if resp != nil {
		sendLog.MessageID = resp.MessageId
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

	cache, err := redis.NewRedisStore()
	if err != nil {
		logrus.WithError(err).Fatal("Redis: can't establish connection")
	}

	config := nsq.NewConfig()

	consumer, err := nsq.NewConsumer(entities.SenderTopic, entities.SenderTopic, config)
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
			storage: s,
			cache:   cache,
		},
		20,
	)

	addr := fmt.Sprintf("%s:%s", os.Getenv("NSQLOOKUPD_HOST"), os.Getenv("NSQLOOKUPD_PORT"))
	nsqlds := []string{addr}

	logrus.Infoln("Connecting to NSQlookup...")
	if err := consumer.ConnectToNSQLookupds(nsqlds); err != nil {
		logrus.WithError(err).Fatal("Nsqlookup: can't establish connection")
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

func newSesClient(keys entities.SesKeys) (emails.Sender, error) {
	if keys.AccessKey == "" || keys.SecretKey == "" || keys.Region == "" {
		return nil, ErrInvalidSesKeys
	}

	client, err := emails.NewSesSender(keys.AccessKey, keys.SecretKey, keys.Region)
	if err != nil {
		return nil, fmt.Errorf("new ses sender: %w", err)
	}

	return client, nil
}

func sendEmail(client emails.Sender, msg entities.SenderTopicParams) (*ses.SendEmailOutput, error) {
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
