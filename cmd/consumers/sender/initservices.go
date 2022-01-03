package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/wire"
	awssqs "github.com/mailbadger/app/sqs"
)

//nolint
var svcSet = wire.NewSet(
	initAwsConfig,
	awssqs.NewClient,
	wire.Bind(new(awssqs.SQSSendReceiveMessageAPI), new(*sqs.Client)),
	awssqs.GetSendEmailQueueURL,
	newQueueURL,
	awssqs.NewConsumerFrom,
)

func initAwsConfig(ctx context.Context) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx)
}

func newQueueURL(url awssqs.SendEmailQueueURL) awssqs.QueueURL {
	return awssqs.QueueURL(url)
}
