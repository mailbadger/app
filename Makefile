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

build: build_static

build_static:
	go install github.com/news-maily/app
	mkdir -p release
	cp $(GOPATH)/bin/api release/

image:
	docker build -t news-maily/api:latest .