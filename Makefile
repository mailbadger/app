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
	mkdir -p release
	go build -o release/app .
	go build -o release/bulksender ./consumers/bulksender
	go build -o release/campaigner ./consumers/campaigner

build_web:
	cd dashboard; yarn && yarn build
	mkdir -p release
	cp -r dashboard/build release/dashboard

image:
	docker build -t news-maily/app:latest .
