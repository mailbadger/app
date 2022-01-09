package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/google/wire"

	"github.com/mailbadger/app/emails"
	boundarysvc "github.com/mailbadger/app/services/boundaries"
	"github.com/mailbadger/app/services/campaigns/scheduler"
	"github.com/mailbadger/app/services/exporters"
	reportsvc "github.com/mailbadger/app/services/reports"
	subscrsvc "github.com/mailbadger/app/services/subscribers"
	templatesvc "github.com/mailbadger/app/services/templates"
	"github.com/mailbadger/app/session"
	awssqs "github.com/mailbadger/app/sqs"
	awss3 "github.com/mailbadger/app/storage/s3"
)

//nolint
var svcSet = wire.NewSet(
	initAwsConfig,
	session.From,
	awssqs.NewClient,
	awss3.NewClient,
	emails.NewSesSender,
	wire.Bind(new(s3iface.S3API), new(*s3.S3)),
	awssqs.GetCampaignerQueueURL,
	wire.Bind(new(awssqs.SendReceiveMessageAPI), new(*sqs.Client)),
	awssqs.NewPublisher,
	wire.Bind(new(awssqs.PublisherAPI), new(awssqs.Publisher)),
	scheduler.New,
	templatesvc.From,
	boundarysvc.New,
	subscrsvc.New,
	exporters.NewSubscribersExporter,
	wire.Bind(new(exporters.Exporter), new(*exporters.SubscribersExporter)),
	reportsvc.New,
)

func initAwsConfig(ctx context.Context) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx)
}
