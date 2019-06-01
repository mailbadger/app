# Builder image
FROM golang:1.12 as go-build

ENV GO111MODULE=on

WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go build -tags json1 github.com/mattn/go-sqlite3

COPY . .

RUN make gen
RUN go build -o /go/bin/app .
RUN go build -o /go/bin/consumers/bulksender ./consumers/bulksender
RUN go build -o /go/bin/consumers/campaigner ./consumers/campaigner

FROM node:8-alpine as node-build

WORKDIR /www/app

COPY dashboard .

RUN yarn build

# Copy into base image
FROM gcr.io/distroless/base

# USER nobody:nobody

ENV APP_DIR=/www/app

COPY --from=go-build /go/bin/app /
COPY --from=go-build /go/bin/consumers /consumers
COPY --from=node-build /www/app/build /www/app/
