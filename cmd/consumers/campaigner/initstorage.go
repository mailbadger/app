package main

import (
	"github.com/google/wire"
	"github.com/mailbadger/app/storage"
)

//nolint
var storeSet = wire.NewSet(storage.New, storage.From)
