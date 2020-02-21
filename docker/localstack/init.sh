#!/usr/bin/env bash
set -euo pipefail

# This script is ran on buildtime to set up the S3 buckets with the files we need

trap "kill %1" EXIT
SERVICES=s3 ./bin/localstack start --host &
while ! curl -fs -o /dev/null "http://localhost:4572"; do
  sleep 10
done

pushd "./data" > /dev/null
find . ! -path . -type d -maxdepth 1 | sed 's/.\//s3:\/\//' | xargs -n1 -I% awslocal s3 mb %
find . -type f | sed 's/.\///' | xargs -n1 -I% awslocal s3 mv % s3://%
popd > /dev/null
