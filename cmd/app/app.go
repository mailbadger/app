//go:build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/mailbadger/app/config"
	"github.com/mailbadger/app/server"
)

type app struct {
	srv *server.Server
}

func newApp(srv *server.Server) app {
	return app{srv}
}

func initApp(conf config.Config) (app, error) {
	wire.Build(storeSet, serverSet, svcSet, newApp)
	return app{}, nil
}
