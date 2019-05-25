#!/usr/bin/env bash

set -euxo pipefail

eval $(egrep -v '^#' .env | xargs) go run newsmaily.go