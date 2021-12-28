package main

import (
	"github.com/google/wire"
	"github.com/mailbadger/app/session"
)

//nolint
var svcSet = wire.NewSet(session.From)
