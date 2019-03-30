# Builder image
FROM golang:1.12 as build

ENV GO111MODULE=on

WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN make gen
RUN go build -o /go/bin/app .
RUN go build -o /go/bin/consumers/bulksender ./consumers/bulksender
RUN go build -o /go/bin/consumers/campaigner ./consumers/campaigner

# Copy into base image
FROM gcr.io/distroless/base
COPY --from=build /go/bin/app /
COPY --from=build /go/bin/consumers /consumers
