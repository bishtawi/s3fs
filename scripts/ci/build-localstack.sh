#!/usr/bin/env bash
set -euo pipefail

if git diff-tree --no-commit-id --name-only -m -r "$GITHUB_SHA" | grep -E "docker/localstack/"; then
    DOCKER_TAG="latest" docker-compose build

    IMAGE="docker.pkg.github.com/bishtawi/s3fs/localstack"
    docker tag "$IMAGE:latest" "$IMAGE:$GITHUB_SHA"

    DOCKER_TAG="$GITHUB_SHA" docker-compose push
    DOCKER_TAG="latest" docker-compose push
fi
