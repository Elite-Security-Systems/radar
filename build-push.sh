#!/bin/bash
set -e

# Get version from git tag or default to dev
VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "dev")
BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ')
COMMIT=$(git rev-parse --short HEAD)

echo "Building RADAR Docker image..."
echo "Version: $VERSION"
echo "Build Date: $BUILD_DATE"
echo "Commit: $COMMIT"

# Build the docker image
docker build \
  --build-arg VERSION=$VERSION \
  --build-arg BUILD_DATE=$BUILD_DATE \
  --build-arg COMMIT=$COMMIT \
  -t elitesecuritysystems/radar:latest \
  -t elitesecuritysystems/radar:$VERSION \
  .

echo "Build complete!"

# Check if we should push the images
if [ "$1" = "--push" ]; then
  echo "Pushing images to Docker Hub..."
  docker push elitesecuritysystems/radar:latest
  docker push elitesecuritysystems/radar:$VERSION
  echo "Push complete!"
fi

echo ""
echo "To use the image, run:"
echo "docker run elitesecuritysystems/radar -domain example.com"
