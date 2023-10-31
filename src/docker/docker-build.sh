#!/bin/bash

# Build Go binary
GOOS=linux GOARCH=amd64 go build -o "institute-person-api" src/main.go
if [ $? -ne 0 ]; then
    echo "Go build failed"
    exit 1
fi

# Get branch and patch level
BRANCH=$(git branch --show-current)
if [ $? -ne 0 ]; then
    echo "Failed to get git branch"
    exit 1
fi

PATCH=$(git rev-parse $BRANCH)
if [ $? -ne 0 ]; then
    echo "Failed to get git commit hash"
    exit 1
fi

# Create PATCH_LEVEL file
echo $BRANCH.$PATCH > PATCH_LEVEL

# Build Docker image
docker build . --tag ghcr.io/agile-learning-institute/institute-person-api:latest
if [ $? -ne 0 ]; then
    echo "Docker build failed"
    exit 1
fi

# Push Docker image
#docker push ghcr.io/agile-learning-institute/institute-person-api:latest
#if [ $? -ne 0 ]; then
#    echo "Docker push failed"
#    exit 1
#fi
#
#docker push ghcr.io/agile-learning-institute/institute-person-api:$BRANCH.$PATCH
#if [ $? -ne 0 ]; then
#    echo "Docker push failed"
#    exit 1
#fi
