#!/usr/bin/env bash

set -euxo pipefail

export $(egrep -v '^#' .env.local | xargs)

go run cmd/sumsubscribermetrics/main.go -start_date="2021-05-01" -end_date="2021-05-20"

#-date="2021-05-02"
#-start_date="2021-05-01" -end_date="2021-05-20"
