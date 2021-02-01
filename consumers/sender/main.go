package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"

	"github.com/mailbadger/app/consumers"
	"github.com/mailbadger/app/emails"
	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/utils"
)

const (
	CharSet = "UTF-8"
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

	msg := new(entities.SenderTopicParams)
	err := json.Unmarshal(m.Body, msg)
	if err != nil {
		logrus.WithField("body", string(m.Body)).
			WithError(err).Error("Malformed JSON message.")
		return nil
	}

	if msg.SesKeys == nil {
		logrus.WithFields(logrus.Fields{
			"uuid":          msg.UUID,
			"user_id":       msg.UserID,
			"campaign_id":   msg.CampaignID,
			"subscriber_id": msg.SubscriberID,
		}).Error("SES Keys are nil.")
		return nil
	}

	client, err := emails.NewSesSender(msg.SesKeys.AccessKey, msg.SesKeys.SecretKey, msg.SesKeys.Region)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"uuid":          msg.UUID,
			"user_id":       msg.UserID,
			"campaign_id":   msg.CampaignID,
			"subscriber_id": msg.SubscriberID,
		}).WithError(err).Error("Unable to create SES sender")
		return nil
	}

	_, err = h.s.GetSendLogByUUID(msg.UUID)
	if err == nil {
		logrus.WithFields(logrus.Fields{
			"uuid":          msg.UUID,
			"user_id":       msg.UserID,
			"campaign_id":   msg.CampaignID,
			"subscriber_id": msg.SubscriberID,
		}).Warn("Email already sent")
		return nil
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

	resp, err := client.SendEmail(input)
	if err != nil {
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
		return nil
	}

	err = h.s.CreateSendLog(&entities.SendLog{
		UUID:         msg.UUID,
		MessageID:    resp.MessageId,
		UserID:       msg.UserID,
		CampaignID:   msg.CampaignID,
		SubscriberID: msg.SubscriberID,
		Status:       entities.StatusDone,
		Description:  resp.GoString(),
	})
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
