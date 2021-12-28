package main

import (
	"context"
	"os/signal"
	"syscall"

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

	app, err := initApp(conf)
	if err != nil {
		logrus.WithError(err).Fatalln("unable to initialize app")
	}

	g := new(errgroup.Group)
	g.Go(func() error {
		return app.srv.ListenAndServe(ctx)
	})

	g.Go(func() error {
		//	cron := metric.NewCron()
		//	return cron.Start(ctx, 5*time.Second)
		return nil
	})

	if err := g.Wait(); err != nil {
		logrus.WithError(err).Error("app terminated")
	}
}
