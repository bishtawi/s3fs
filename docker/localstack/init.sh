#!/usr/bin/env bash
set -euo pipefail

# Initialize the S3 buckets

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

if [ "${1:-}" == "true" ]; then
  trap "kill %1" EXIT
  SERVICES=s3 ./bin/localstack start --host &
  while ! curl -fs -o /dev/null "http://localhost:4572"; do
    sleep 10
  done
fi

pushd "$DIR/data" > /dev/null
find . ! -path . -type d -maxdepth 1 | sed 's/.\//s3:\/\//' | xargs -n1 -I% awslocal s3 mb %
find . -type f | sed 's/.\///' | xargs -n1 -I% awslocal s3 cp % s3://%
popd > /dev/null
