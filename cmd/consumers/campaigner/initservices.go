package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/google/wire"
	"github.com/mailbadger/app/services/campaigns"
	"github.com/mailbadger/app/services/templates"
	awssqs "github.com/mailbadger/app/sqs"
	awss3 "github.com/mailbadger/app/storage/s3"
)

//nolint
var svcSet = wire.NewSet(
	initAwsConfig,
	awss3.NewClient,
	awssqs.NewClient,
	wire.Bind(new(awssqs.SendReceiveMessageAPI), new(*sqs.Client)),
	wire.Bind(new(s3iface.S3API), new(*s3.S3)),
	awssqs.GetCampaignerQueueURL,
	awssqs.GetSendEmailQueueURL,
	newQueueURL,
	awssqs.NewPublisher,
	awssqs.NewConsumerFrom,
	templates.New,
	campaigns.New,
)

func initAwsConfig(ctx context.Context) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx)
}

func newQueueURL(url awssqs.CampaignerQueueURL) awssqs.QueueURL {
	return awssqs.QueueURL(url)
}
