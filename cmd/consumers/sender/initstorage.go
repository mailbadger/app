package main

import (
	"github.com/google/wire"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/storage/redis"
)

//nolint
var storeSet = wire.NewSet(
	storage.New,
	storage.From,
	redis.NewStoreFrom,
	wire.Bind(new(redis.Store), new(*redis.RedisStore)),
)
