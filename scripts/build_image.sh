#!/bin/bash

if [ -n "$1" ]; then
  APP_NAME="$1"
else
  echo "Error: APP_NAME is required"
  exit 1
fi

if [ -n "$2" ]; then
  VERSION="$2"
else
  VERSION=$(git describe --tags --always --dirty)
fi

BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ')

# Build the Docker image
echo "Building Docker image ${APP_NAME}:${VERSION}..."

docker build -f Dockerfile \
    --build-arg VERSION="${VERSION}" \
    --build-arg BUILD_DATE="${BUILD_DATE}" \
    -t "${APP_NAME}:${VERSION}" \
    -t "${APP_NAME}:latest" \
    . || exit 1

echo "Docker image built successfully."
echo "Tags:"
echo "  - ${APP_NAME}:${VERSION}"
echo "  - ${APP_NAME}:latest"
