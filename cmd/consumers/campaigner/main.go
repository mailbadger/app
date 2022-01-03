package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/mailbadger/app/config"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	conf, err := config.FromEnv()
	if err != nil {
		logrus.WithError(err).Fatalln("unable to read config from env")
	}

	initMode(conf.Mode)
	initLogger(conf.Logging)

	app, err := initApp(ctx, conf)
	if err != nil {
		logrus.WithError(err).Fatalln("unable to initialize app")
	}

	fn := func(ctx context.Context, m types.Message) func() error {
		return func() error {
			err = app.handler.HandleMessage(ctx, m)
			if err != nil {
				return err
			}
			return app.handler.DeleteMessage(ctx, m)
		}
	}

	g := new(errgroup.Group)
	for m := range app.consumer.PollSQS(ctx) {
		g.Go(fn(ctx, m))
	}

	if err := g.Wait(); err != nil {
		logrus.WithError(err).Error("received an error when handling a message")
	}
}
