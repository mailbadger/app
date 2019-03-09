# Start by building the application.
FROM golang:1.12 as build

WORKDIR /go/src/app
COPY . .

RUN go-wrapper download
RUN make gen
RUN go-wrapper install


# Now copy it into our base image.
FROM gcr.io/distroless/base
COPY --from=build /go/bin/app /
CMD ["/app"]
