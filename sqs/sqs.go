package sqs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/gin-gonic/gin"
	"github.com/mailbadger/app/config"
	"github.com/mailbadger/app/logger"
)

const (
	// CampaignerTopic is the topic used by the campaigner consumer.
	CampaignerTopic = "SendCampaign"
	// SenderTopic is the topic used by the sender consumer.
	SenderTopic = "SendEmail"
)

// QueueURL is a pointer to a URL string, used by the SQS client.
type QueueURL *string

// CampaignerQueueURL represents the queue url of the SendCampaign queue.
type CampaignerQueueURL QueueURL

// CampaignerQueueURL represents the queue url of the SendEmail queue.
type SendEmailQueueURL QueueURL

// SQSSendReceiveMessageAPI defines the interface for the GetQueueUrl function.
// We use this interface to test the function using a mocked service.
type SQSSendReceiveMessageAPI interface {
	GetQueueUrl(ctx context.Context,
		params *sqs.GetQueueUrlInput,
		optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)

	ReceiveMessage(ctx context.Context,
		params *sqs.ReceiveMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)

	SendMessage(ctx context.Context,
		params *sqs.SendMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

type Consumer struct {
	queueURL          QueueURL
	visibilityTimeout int32
	maxNumOfMessages  int32
	waitTimeout       int32
	api               SQSSendReceiveMessageAPI
}

type Publisher struct {
	api SQSSendReceiveMessageAPI
}

func NewClient(cfg aws.Config) *sqs.Client {
	return sqs.NewFromConfig(cfg)
}

func NewConsumerFrom(
	conf config.Config,
	queueURL QueueURL,
	api SQSSendReceiveMessageAPI,
) Consumer {
	return NewConsumer(
		queueURL,
		conf.Consumer.Timeout,
		conf.Consumer.WaitTimeout,
		conf.Consumer.MaxInFlightMsgs,
		api,
	)
}

func NewConsumer(
	queueURL QueueURL,
	visibilityTimeout,
	maxNumOfMessages,
	waitTimeout int32,
	api SQSSendReceiveMessageAPI,
) Consumer {
	return Consumer{
		queueURL:          queueURL,
		visibilityTimeout: visibilityTimeout,
		maxNumOfMessages:  maxNumOfMessages,
		waitTimeout:       waitTimeout,
		api:               api,
	}
}

func NewPublisher(api SQSSendReceiveMessageAPI) Publisher {
	return Publisher{
		api: api,
	}
}

func (c Consumer) PollSQS(ctx context.Context) <-chan types.Message {
	msgs := make(chan types.Message)
	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.From(ctx).Info("sqs consumer: polling canceled...")
				close(msgs)
				return
			default:
				gMInput := &sqs.ReceiveMessageInput{
					MessageAttributeNames: []string{
						string(types.QueueAttributeNameAll),
					},
					AttributeNames: []types.QueueAttributeName{
						types.QueueAttributeName(types.MessageSystemAttributeNameSentTimestamp),
					},
					QueueUrl:            c.queueURL,
					MaxNumberOfMessages: c.maxNumOfMessages,
					VisibilityTimeout:   c.visibilityTimeout,
					WaitTimeSeconds:     c.waitTimeout,
				}

				msgResult, err := getMessages(ctx, c.api, gMInput)
				if err != nil {
					logger.From(ctx).WithError(err).Error("sqs consumer: unable to get messages, aborting...")
					close(msgs)
					return
				}

				for _, m := range msgResult.Messages {
					msgs <- m
				}
			}
		}
	}()

	return msgs
}

// GetMessages gets the most recent message from an Amazon SQS queue.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a ReceiveMessageOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to ReceiveMessage.
func getMessages(
	ctx context.Context,
	api SQSSendReceiveMessageAPI,
	input *sqs.ReceiveMessageInput,
) (*sqs.ReceiveMessageOutput, error) {
	return api.ReceiveMessage(ctx, input)
}

func (p Publisher) SendMessage(ctx context.Context, queueUrl *string, body []byte) error {
	b := string(body)
	_, err := p.api.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: &b,
		QueueUrl:    queueUrl,
	})

	return err
}

func (p Publisher) GetQueueURL(ctx context.Context, queueName *string) (*string, error) {
	out, err := p.api.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: queueName})
	if err != nil {
		return nil, err
	}
	return out.QueueUrl, nil
}

const key = "publisher"

// SetPublisherToContext sets the producer to the context
func SetPublisherToContext(ctx *gin.Context, pub Publisher) {
	ctx.Set(key, pub)
}

// GetPublisherFromContext returns the Producer associated with the context
func GetPublisherFromContext(ctx context.Context) Publisher {
	return ctx.Value(key).(Publisher)
}

func SendMessage(ctx context.Context, queueUrl *string, body []byte) error {
	return GetPublisherFromContext(ctx).SendMessage(ctx, queueUrl, body)
}

func GetQueueURL(ctx context.Context, queueName *string) (*string, error) {
	return GetPublisherFromContext(ctx).GetQueueURL(ctx, queueName)
}

func GetCampaignerQueueURL(ctx context.Context, api SQSSendReceiveMessageAPI) (CampaignerQueueURL, error) {
	queueStr := CampaignerTopic
	gQInput := &sqs.GetQueueUrlInput{
		QueueName: &queueStr,
	}
	// Get URL of queue
	urlResult, err := api.GetQueueUrl(ctx, gQInput)
	if err != nil {
		return nil, err
	}
	return urlResult.QueueUrl, nil
}

func GetSendEmailQueueURL(ctx context.Context, api SQSSendReceiveMessageAPI) (SendEmailQueueURL, error) {
	queueStr := SenderTopic
	gQInput := &sqs.GetQueueUrlInput{
		QueueName: &queueStr,
	}
	// Get URL of queue
	urlResult, err := api.GetQueueUrl(ctx, gQInput)
	if err != nil {
		return nil, err
	}
	return urlResult.QueueUrl, nil
}
