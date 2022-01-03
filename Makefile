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
	statik -ns=migrations -src=./storage/migrations/$(driver) -f

test: 
	go test -cover $(PACKAGES)

build: build_api

build_api:
	mkdir -p bin
	go build -o bin/app ./cmd/app
	go build -o bin/sender ./cmd/consumers/sender
	go build -o bin/campaigner ./cmd/consumers/campaigner

build_static:
	cd dashboard; rm -rf build && yarn && yarn build

image:
	docker build -t mailbadger/app:latest .

run_mailbadger:
	./scripts/run-mailbadger.sh

run_campaigner:
	./scripts/run-campaigner.sh

run_sender:
	./scripts/run-sender.sh

install_fixtures:
	./scripts/install-fixtures.sh

process_events:
	./scripts/process-events.sh
