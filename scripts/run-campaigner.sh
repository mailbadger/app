#!/usr/bin/env bash

set -euxo pipefail

export $(egrep -v '^#' .env.local | xargs)

go run -trace=true consumers/campaigner/main.go 2> trace.out
