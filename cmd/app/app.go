//go:build wireinject

package main

import (
	"context"
	"github.com/google/wire"

	"github.com/mailbadger/app/config"
	"github.com/mailbadger/app/server"
	"github.com/mailbadger/app/services/campaigns/scheduler"
	"github.com/mailbadger/app/services/subscribers/metric"
)

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

func initApp(ctx context.Context, conf config.Config) (app, error) {
	wire.Build(storeSet, serverSet, svcSet, newApp)
	return app{}, nil
}
