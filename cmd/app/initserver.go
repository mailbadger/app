package main

import (
	"github.com/google/wire"
	"github.com/mailbadger/app/opa"
	"github.com/mailbadger/app/routes"
	"github.com/mailbadger/app/server"
)

//nolint
var serverSet = wire.NewSet(opa.NewCompiler, routes.From, server.From)
