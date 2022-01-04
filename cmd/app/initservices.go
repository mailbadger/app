package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/wire"
	"github.com/mailbadger/app/services/campaigns/scheduler"
	"github.com/mailbadger/app/services/subscribers/metric"
	"github.com/mailbadger/app/session"
	awssqs "github.com/mailbadger/app/sqs"
	"github.com/mailbadger/app/storage/s3"
)

//nolint
var svcSet = wire.NewSet(
	initAwsConfig,
	session.From,
	awssqs.NewClient,
	s3.NewClient,
	awssqs.GetCampaignerQueueURL,
	wire.Bind(new(awssqs.SendReceiveMessageAPI), new(*sqs.Client)),
	awssqs.NewPublisher,
	metric.NewCron,
	scheduler.New,
)

func initAwsConfig(ctx context.Context) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx)
}
