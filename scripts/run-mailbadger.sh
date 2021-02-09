#!/usr/bin/env bash

set -euxo pipefail

make gen

export $(egrep -v '^#' .env.local | xargs)

go run mailbadger.go
