package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/news-maily/api/consumers"

	"github.com/news-maily/api/emails"
	"github.com/news-maily/api/entities"
	"github.com/sirupsen/logrus"

	"github.com/news-maily/api/storage"
	"github.com/nsqio/go-nsq"
)

type MessageHandler struct {
	s storage.Storage
}

// HandleMessage is the only requirement needed to fulfill the
// nsq.Handler interface.
func (h *MessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return errors.New("body is blank")
	}

	msg := new(entities.BulkSendMessage)

	err := json.Unmarshal(m.Body, msg)
	if err != nil {
		return err
	}

	if msg.SesKeys == nil {
		return errors.New("SES Keys are nil")
	}

	client, err := emails.NewSesSender(msg.SesKeys.AccessKey, msg.SesKeys.SecretKey, msg.SesKeys.Region)
	if err != nil {
		return err
	}

	res, err := client.SendBulkTemplatedEmail(msg.Input)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"campaign_id":   msg.CampaignID,
			"template_name": msg.Input.Template,
			"user_id":       msg.UserID,
		}).Errorln(err.Error())

		return nil //we won't re-throw the message here.
	}

	for _, s := range res.Status {
		err := h.s.CreateSendBulkLog(&entities.SendBulkLog{
			UserID:     msg.UserID,
			CampaignID: msg.CampaignID,
			MessageID:  *s.MessageId,
			Status:     *s.Status,
			Error:      s.Error,
		})
		logrus.WithFields(logrus.Fields{
			"campaign_id":      msg.CampaignID,
			"user_id":          msg.UserID,
			"send_bulk_status": s.GoString(),
		}).Errorln(err.Error())
	}

	return nil
}

func main() {
	lvl, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logrus.Panic(err)
	}

	logrus.SetLevel(lvl)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	s := storage.New(os.Getenv("DATABASE_DRIVER"), os.Getenv("DATABASE_CONFIG"))

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
