#!/usr/bin/env bash

set -euxo pipefail

make gen

export $(egrep -v '^#' .env.local | xargs)

go run consumers/sender/main.go
