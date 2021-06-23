package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/mailbadger/app/mode"
	"github.com/mailbadger/app/routes"
	"github.com/mailbadger/app/server"
)

func init() {
	mode.SetModeFromEnv()
	lvl, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		lvl = logrus.InfoLevel
	}

	logrus.SetLevel(lvl)
	logrus.SetOutput(os.Stdout)
	if mode.IsProd() {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		gin.SetMode(gin.ReleaseMode)
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	handler := routes.New()

	var addr = os.Getenv("PORT")
	if addr == "" {
		addr = "8080"
	}
	srv := server.New(
		":"+addr,
		server.WithHandler(handler),
		server.WithTLS(os.Getenv("CERT_FILE"), os.Getenv("KEY_FILE")),
	)

	var g errgroup.Group
	g.Go(func() error {
		logrus.Infoln("Starting the server...")
		return srv.ListenAndServe(ctx)
	})

	if err := g.Wait(); err != nil {
		logrus.WithError(err).Fatalln("Server terminated...")
	}
}
