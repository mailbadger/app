package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"

	"github.com/mailbadger/app/emails"
	"github.com/mailbadger/app/entities"
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

type handler struct {
	storage   storage.Storage
	cache     redis.Store
	sqsclient *sqs.Client
	queueURL  awssqs.SendEmailQueueURL
}

func newHandler(
	storage storage.Storage,
	cache redis.Store,
	queueURL awssqs.SendEmailQueueURL,
) *handler {
	return &handler{
		storage:  storage,
		cache:    cache,
		queueURL: queueURL,
	}
}

// HandleMessage is the only requirement needed to fulfill the
// nsq.Handler interface.
func (h *handler) HandleMessage(ctx context.Context, m types.Message) (err error) {
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

	logEntry.Info("Received message, processing..")

	// check if the message is processing (if the uuid exists in redis that means it is in progress)
	exist, err := h.cache.Exists(ctx, cacheKey)
	if err != nil {
		logEntry.WithError(err).Error("Unable to check message existence")
		return err
	}

	if exist {
		logEntry.WithError(err).Info("Message already processed")
		return nil
	}

	if err := h.cache.Set(ctx, cacheKey, []byte("1"), cacheDuration); err != nil {
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
				rerr := h.cache.Delete(ctx, cacheKey)
				if rerr != nil {
					logEntry.WithError(rerr).Error("Unable to delete cached id")
				}
				return err
			case sns.ErrCodeInternalErrorException:
				logEntry.WithError(aerr).Error("Unable to send templated email. The request processing has failed because of an unknown error, exception, or failure.")
				rerr := h.cache.Delete(ctx, cacheKey)
				if rerr != nil {
					logEntry.WithError(rerr).Error("Unable to delete cached id")
				}
				return err
			default:
				logEntry.WithError(aerr).Error("Unable to send templated email. Unknown status.")
				rerr := h.cache.Delete(ctx, cacheKey)
				if rerr != nil {
					logEntry.WithError(rerr).Error("Unable to delete cached id")
				}
				return err
			}
		} else {
			logEntry.WithError(err).Error("Unable to send templated email.")
			rerr := h.cache.Delete(ctx, cacheKey)
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

func (h *handler) DeleteMessage(ctx context.Context, m types.Message) error {
	_, err := h.sqsclient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      h.queueURL,
		ReceiptHandle: m.ReceiptHandle,
	})
	return err
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
	return prefix + hex.EncodeToString(k)
}
