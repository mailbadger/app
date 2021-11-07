#!/usr/bin/env bash

set -euxo pipefail

export $(egrep -v '^#' .env.local | xargs)
trap 'kill 0' SIGINT; go run ./cmd/app/main.go & \
  go run ./cmd/consumers/campaigner/main.go & \
  go run ./cmd/consumers/sender/main.go
