// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"github.com/mailbadger/app/config"
	"github.com/mailbadger/app/opa"
	"github.com/mailbadger/app/routes"
	"github.com/mailbadger/app/server"
	"github.com/mailbadger/app/services/campaigns/scheduler"
	"github.com/mailbadger/app/services/subscribers/metric"
	"github.com/mailbadger/app/session"
	"github.com/mailbadger/app/sqs"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/storage/s3"
)

// Injectors from app.go:

func initApp(ctx context.Context, conf config.Config) (app, error) {
	db := storage.New(conf)
	storageStorage := storage.From(db)
	sessionSession := session.From(storageStorage, conf)
	compiler, err := opa.NewCompiler()
	if err != nil {
		return app{}, err
	}
	awsConfig, err := initAwsConfig(ctx)
	if err != nil {
		return app{}, err
	}
	client := sqs.NewClient(awsConfig)
	publisher := sqs.NewPublisher(client)
	s3S3, err := s3.NewClient()
	if err != nil {
		return app{}, err
	}
	api := routes.From(sessionSession, storageStorage, compiler, publisher, s3S3, conf)
	serverServer := server.From(api, conf)
	cron := metric.NewCron(storageStorage)
	campaignerQueueURL, err := sqs.GetCampaignerQueueURL(ctx, client)
	if err != nil {
		return app{}, err
	}
	schedulerScheduler := scheduler.New(storageStorage, publisher, campaignerQueueURL)
	mainApp := newApp(serverServer, cron, schedulerScheduler)
	return mainApp, nil
}

// app.go:

type app struct {
	srv           *server.Server
	subscrmetrics *metric.Cron
	campaignsched *scheduler.Scheduler
}

func newApp(
	srv *server.Server,
	subscrmetrics *metric.Cron,
	campaignsched *scheduler.Scheduler,
) app {
	return app{
		srv:           srv,
		subscrmetrics: subscrmetrics,
		campaignsched: campaignsched,
	}
}
