package main

import (
	"os"

	"github.com/mailbadger/app/config"
	"github.com/mailbadger/app/mode"
	"github.com/sirupsen/logrus"
)

// nolint
func initLogger(logConf config.Logging) {
	lvl, err := logrus.ParseLevel(logConf.Level)
	if err != nil {
		lvl = logrus.InfoLevel
	}

	logrus.SetLevel(lvl)
	logrus.SetOutput(os.Stdout)
	if mode.IsProd() {
		logrus.SetFormatter(&logrus.JSONFormatter{
			PrettyPrint: logConf.Pretty,
		})
	}
}
