package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/mailbadger/app/emails"
	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/mode"
	awssqs "github.com/mailbadger/app/sqs"
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

	msg := new(entities.SenderTopicParams)
	err = json.Unmarshal([]byte(*m.Body), msg)
	if err != nil {
		logrus.WithField("body", string(*m.Body)).
			WithError(err).Error("Malformed JSON message.")
		return err
	}

	cacheKey := genCacheKey(cachePrefix, fmt.Sprintf("%s_%d", msg.EventID, msg.SubscriberID))
	logEntry := logrus.WithFields(logrus.Fields{
		"event_id":      msg.EventID,
		"user_id":       msg.UserID,
		"campaign_id":   msg.CampaignID,
		"subscriber_id": msg.SubscriberID,
		"cache_key":     cacheKey,
	})

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

	if err := h.cache.Set(cacheKey, []byte("1"), cacheDuration); err != nil {
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
	timeout := flag.Int("t", 300, "How long, in seconds, that the message is hidden from others")
	maxInFlightMsgs := flag.Int("m", 100, "Max number of messages to be received from SQS simultaneously")
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
	consumer := awssqs.NewConsumer(
		*queueURL,
		int32(*timeout),
		int32(*maxInFlightMsgs),
		int32(*waitTimeout),
		client,
	)

	driver := os.Getenv("DATABASE_DRIVER")
	conf := storage.MakeConfigFromEnv(driver)
	s := storage.New(driver, conf)

	cache, err := redis.NewRedisStore()
	if err != nil {
		logrus.WithError(err).Fatal("Redis: can't establish connection")
	}

	g, ctx := errgroup.WithContext(ctx)
	handler := &MessageHandler{
		storage:   s,
		cache:     cache,
		sqsclient: client,
		queueURL:  queueURL,
	}

	fn := func(m types.Message) func() error {
		return func() error {
			err := handler.HandleMessage(ctx, m)
			if err != nil {
				// on error we must not delete the message. We want other
				// consumers to try and process the message again.
				return err
			}
			return handler.DeleteMessage(ctx, m)
		}
	}

	// Poll for messages
	messages := consumer.PollSQS(ctx)

	for m := range messages {
		g.Go(fn(m))
	}

	if err := g.Wait(); err != nil {
		logrus.WithError(err).Error("received an error when handling a message")
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

func genCacheKey(prefix string, key string) string {
	h := sha256.New()
	h.Write([]byte(key))
	k := h.Sum(nil)
	return prefix + string(k)
}
