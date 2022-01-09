package main

import (
	"github.com/google/wire"
	"github.com/mailbadger/app/session"
	"github.com/mailbadger/app/storage"
)

//nolint
var storeSet = wire.NewSet(storage.New, storage.From, wire.Bind(new(session.Store), new(storage.Storage)))
