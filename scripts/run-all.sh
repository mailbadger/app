#!/usr/bin/env bash

set -euxo pipefail

export $(egrep -v '^#' .env.local | xargs)
trap 'kill 0' SIGINT; go run ./cmd/app/... & \
  go run ./cmd/consumers/campaigner/... & \
  go run ./cmd/consumers/sender/...
