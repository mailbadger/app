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

# Copy into base image
FROM gcr.io/distroless/base
COPY --from=build /go/bin/app /
CMD ["/app"]
