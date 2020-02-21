#!/usr/bin/env bash
set -euo pipefail

if [ "$#" -lt 1 ]; then
  echo 's3local.sh wraps the aws s3 cli to point to localstack s3 service. Available commands:'
  aws s3 error
fi

ENDPOINT_URL="http://localhost:4572"
export AWS_ACCESS_KEY_ID=s3local.sh AWS_SECRET_ACCESS_KEY=localstack

aws "--endpoint-url=$ENDPOINT_URL" s3 "${@}"
