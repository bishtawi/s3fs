#!/usr/bin/env bash
set -euo pipefail

if git diff-tree --no-commit-id --name-only -m -r "$GITHUB_SHA" | grep -E "docker/localstack/"; then
    DOCKER_TAG="latest" docker-compose build

    IMAGES=$(docker-compose config | grep 'image: ' | cut -d':' -f 2 | tr -d '"')
    for IMAGE in $IMAGES; do
        docker tag "$IMAGE:latest" "$IMAGE:$GITHUB_SHA"
    done

    DOCKER_TAG="$GITHUB_SHA" docker-compose push
    DOCKER_TAG="latest" docker-compose push
fi
