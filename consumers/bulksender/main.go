package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/news-maily/app/consumers"
	"github.com/news-maily/app/emails"
	"github.com/news-maily/app/entities"
	"github.com/news-maily/app/storage"
	"github.com/news-maily/app/utils"
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
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
		}).Errorln(err.Error())
		return nil
	}

	count, err := h.s.CountLogsByUUID(msg.UUID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id":     msg.UserID,
			"campaign_id": msg.CampaignID,
			"uuid":        msg.UUID,
		}).Errorln(err.Error())
		return nil
	}

	if count > 0 {
		logrus.WithFields(logrus.Fields{
			"user_id":     msg.UserID,
			"campaign_id": msg.CampaignID,
			"uuid":        msg.UUID,
		}).Warn("bulk already sent")
		return nil
	}

	res, err := client.SendBulkTemplatedEmail(msg.Input)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id":     msg.UserID,
			"campaign_id": msg.CampaignID,
			"uuid":        msg.UUID,
		}).Errorln(err.Error())

		// TODO - handle throttle exceptions from SES and re-queue accordingly.
		return nil
	}

	for _, s := range res.Status {
		err := h.s.CreateSendBulkLog(&entities.SendBulkLog{
			UUID:       msg.UUID,
			UserID:     msg.UserID,
			CampaignID: msg.CampaignID,
			MessageID:  *s.MessageId,
			Status:     *s.Status,
		})
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"campaign_id":      msg.CampaignID,
				"user_id":          msg.UserID,
				"send_bulk_status": s.GoString(),
			}).Errorln(err.Error())
		}
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

	consumer, err := nsq.NewConsumer(entities.SendBulkTopic, entities.SendBulkTopic, config)
	if err != nil {
		log.Fatal(err)
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
	if err := consumer.ConnectToNSQLookupds(nsqlds); err != nil {
		log.Fatal(err)
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
			consumer.Stop()
		}
	}
}
