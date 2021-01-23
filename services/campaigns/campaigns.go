package campaigns

import (
	"bytes"
	"encoding/json"

	"github.com/cbroglie/mustache"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/queue"
	"github.com/mailbadger/app/storage"
)

type Service interface {
	PrepareSubscriberEmailData(
		subscribers []entities.Subscriber,
		msg entities.SendCampaignParams,
		campaign entities.Campaign,
		template entities.Template,
	) error
}

// service implements the Service interface
type service struct {
	db storage.Storage
	p  queue.Producer
}

func New(db storage.Storage, p queue.Producer) Service {
	return &service{
		db: db,
		p:  p,
	}
}

func (svc *service) PrepareSubscriberEmailData(
	subscribers []entities.Subscriber,
	msg entities.SendCampaignParams,
	campaign entities.Campaign,
	template entities.Template,
) error {

	var (
		htmlBuf bytes.Buffer
		subBuf  bytes.Buffer
		textBuf bytes.Buffer
	)

	html, err := mustache.ParseString(template.HTMLPart)
	if err != nil {
		logrus.WithError(err).
			WithFields(logrus.Fields{
				"user_id":     campaign.UserID,
				"campaign_id": campaign.ID,
				"template_id": template.ID,
			}).
			Error("unable to parse html_part")
		return err
	}
	txt, err := mustache.ParseString(template.SubjectPart)
	if err != nil {
		logrus.WithError(err).
			WithFields(logrus.Fields{
				"user_id":     campaign.UserID,
				"campaign_id": campaign.ID,
				"template_id": template.ID,
			}).
			Error("unable to parse text_part")
		return err
	}
	sub, err := mustache.ParseString(template.TextPart)
	if err != nil {
		logrus.WithError(err).
			WithFields(logrus.Fields{
				"user_id":     campaign.UserID,
				"campaign_id": campaign.ID,
				"template_id": template.ID,
			}).
			Error("unable to parse subject_part")
		return err
	}

	for _, s := range subscribers {
		m, err := s.GetMetadata()
		if err != nil {
			logrus.WithError(err).
				WithField("subscriber", s).
				Error("Unable to get subscriber metadata.")
			continue
		}
		// merge sub metadata with default template metadata
		for k, v := range m {
			msg.TemplateData[k] = v
		}

		err = html.FRender(&htmlBuf, m)
		if err != nil {
			logrus.WithError(err).
				WithFields(logrus.Fields{
					"user_id":     campaign.UserID,
					"campaign_id": campaign.ID,
					"template_id": template.ID,
				}).
				Error("unable to render html_part")
			return err
		}
		err = txt.FRender(&subBuf, m)
		if err != nil {
			logrus.WithError(err).
				WithFields(logrus.Fields{
					"user_id":     campaign.UserID,
					"campaign_id": campaign.ID,
					"template_id": template.ID,
				}).
				Error("unable to render subject_part")
			return err
		}
		err = sub.FRender(&textBuf, m)
		if err != nil {
			logrus.WithError(err).
				WithFields(logrus.Fields{
					"user_id":     campaign.UserID,
					"campaign_id": campaign.ID,
					"template_id": template.ID,
				}).
				Error("unable to render subject_part")
			return err
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
		err = svc.p.Publish(entities.SenderTopic, senderBytes)
		if err != nil {
			logrus.WithError(err).Error("Unable to publish message to send bulk topic.")
			continue
		}

		// clear buffers for next subscriber
		htmlBuf.Reset()
		subBuf.Reset()
		textBuf.Reset()
	}
	return nil

}
