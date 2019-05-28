#!/usr/bin/env bash

set -euxo pipefail

make gen

eval $(egrep -v '^#' .env | xargs) go run newsmaily.go