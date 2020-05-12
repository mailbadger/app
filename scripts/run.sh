#!/usr/bin/env bash

set -euxo pipefail

make gen

eval $(egrep -v '^#' .env.local | xargs) go run mailbadger.go
