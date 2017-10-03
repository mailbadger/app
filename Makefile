.PHONY: build

ifneq ($(shell uname), Darwin)
	EXTLDFLAGS = -extldflags "-static" $(null)
else
	EXTLDFLAGS =
endif

PACKAGES = $(shell go list ./... | grep -v /vendor/)


gen: gen_migrations

gen_migrations: 
	go generate github.com/FilipNikolovski/news-maily/storage/migrations

deps:
	go get -u github.com/jteeuwen/go-bindata/...
	go get -u github.com/stretchr/testify

test: 
	go test -cover $(PACKAGES)

build: build_static

build_static:
	go install github.com/FilipNikolovski/news-maily
	mkdir -p release
	cp $(GOPATH)/bin/news-maily release/