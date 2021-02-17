.PHONY: build

ifneq ($(shell uname), Darwin)
	EXTLDFLAGS = -extldflags "-static" $(null)
else
	EXTLDFLAGS =
endif

PACKAGES = $(shell go list ./... | grep -v /vendor/)
driver ?= sqlite3


gen: gen_migrations

gen_migrations:
	statik -ns=migrations -src=./storage/migrations/$(driver)

test: 
	go test -cover $(PACKAGES)

build: build_api

build_api:
	mkdir -p bin
	go build -o bin/app .
	go build -o bin/bulksender ./consumers/bulksender
	go build -o bin/campaigner ./consumers/campaigner

build_static:
	cd dashboard; rm -rf build && yarn && yarn build

image:
	docker build -t mailbadger/app:latest .
