.PHONY: build

ifneq ($(shell uname), Darwin)
	EXTLDFLAGS = -extldflags "-static" $(null)
else
	EXTLDFLAGS =
endif

PACKAGES = $(shell go list ./... | grep -v /vendor/)


gen: gen_migrations

gen_migrations:
	go generate github.com/news-maily/app/storage/migrations

test: 
	go test -cover $(PACKAGES)

build: build_api

build_api:
	mkdir -p bin
	go build -o bin/app .
	go build -o bin/bulksender ./consumers/bulksender
	go build -o bin/campaigner ./consumers/campaigner

build_static:
	cd dashboard; yarn && yarn build

image:
	docker build -t news-maily/app:latest .
