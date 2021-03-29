#!/usr/bin/env bash

set -euxo pipefail

export $(egrep -v '^#' .env.local | xargs)

go run consumers/campaigner/main.go
