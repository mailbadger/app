package events

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type EventsClient interface {
	CreateTopic(input *sns.CreateTopicInput) (*sns.CreateTopicOutput, error)
	DeleteTopic(input *sns.DeleteTopicInput) (*sns.DeleteTopicOutput, error)
	Subscribe(input *sns.SubscribeInput) (*sns.SubscribeOutput, error)
}

const SNSTopicName = "MailbadgerEvents"

type eventsClientImpl struct {
	*sns.SNS
}

func NewSNSClient(key, secret, region string) (*sns.SNS, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
	})

	if err != nil {
		return nil, err
	}

	return sns.New(sess), nil
}

func NewEventsClient(key, secret, region string) (EventsClient, error) {
	client, err := NewSNSClient(key, secret, region)
	if err != nil {
		return nil, err
	}

	return &eventsClientImpl{client}, nil
}
