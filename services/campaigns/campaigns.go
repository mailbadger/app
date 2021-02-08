package campaigns

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/cbroglie/mustache"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/queue"
	"github.com/mailbadger/app/storage"
)

type Service interface {
	PrepareSubscriberEmailData(
		s entities.Subscriber,
		uuid string,
		msg entities.CampaignerTopicParams,
		campaignID int64,
		html *mustache.Template,
		sub *mustache.Template,
		text *mustache.Template,
	) (*entities.SenderTopicParams, error)
	PublishSubscriberEmailParams(params *entities.SenderTopicParams) error
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
	s entities.Subscriber,
	uuid string,
	msg entities.CampaignerTopicParams,
	campaignID int64,
	html *mustache.Template,
	sub *mustache.Template,
	text *mustache.Template,
) (*entities.SenderTopicParams, error) {

	var (
		htmlBuf bytes.Buffer
		subBuf  bytes.Buffer
		textBuf bytes.Buffer
	)

	m, err := s.GetMetadata()
	if err != nil {
		return nil, fmt.Errorf("campaign service: prepare email data: get metadata: %w", err)
	}
	// merge sub metadata with default template metadata
	for k, v := range msg.TemplateData {
		if _, ok := m[k]; !ok {
			m[k] = v
		}
	}

	if s.Name != "" {
		m["name"] = s.Name
	}

	url, err := s.GetUnsubscribeURL(msg.UserUUID)
	if err != nil {
		return nil, fmt.Errorf("campaign service: get unsubscribe url: %w", err)
	} else {
		m["unsubscribe_url"] = url
	}

	err = html.FRender(&htmlBuf, m)
	if err != nil {
		return nil, fmt.Errorf("campaign service: prepare email data: render html: %w", err)
	}
	err = sub.FRender(&subBuf, m)
	if err != nil {
		return nil, fmt.Errorf("campaign service: prepare email data: render subject: %w", err)
	}
	err = text.FRender(&textBuf, m)
	if err != nil {
		return nil, fmt.Errorf("campaign service: prepare email data: render text: %w", err)
	}

	sender := entities.SenderTopicParams{
		UUID:                   uuid,
		SubscriberID:           s.ID,
		SubscriberEmail:        s.Email,
		Source:                 msg.Source,
		ConfigurationSetExists: msg.ConfigurationSetExists,
		CampaignID:             campaignID,
		SesKeys:                msg.SesKeys,
		HTMLPart:               htmlBuf.Bytes(),
		SubjectPart:            subBuf.Bytes(),
		TextPart:               textBuf.Bytes(),
		UserUUID:               msg.UserUUID,
		UserID:                 msg.UserID,
	}

	return &sender, nil

}

func (svc *service) PublishSubscriberEmailParams(params *entities.SenderTopicParams) error {
	senderBytes, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("campaign service: publish to sender: marshal params: %w", err)
	}

	// publish the message to the queue
	err = svc.p.Publish(entities.SenderTopic, senderBytes)
	if err != nil {
		return fmt.Errorf("campaign service: publish to sender: %w", err)
	}
	return nil
}
