//go:build wireinject

package main

import (
	"context"

	"github.com/google/wire"
	"github.com/mailbadger/app/config"
	"github.com/mailbadger/app/sqs"
)

type app struct {
	handler  *handler
	consumer sqs.Consumer
}

func newApp(h *handler, c sqs.Consumer) app {
	return app{
		handler:  h,
		consumer: c,
	}
}

func initApp(ctx context.Context, conf config.Config) (app, error) {
	wire.Build(storeSet, svcSet, newHandler, newApp)
	return app{}, nil
}
