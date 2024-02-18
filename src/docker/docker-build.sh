#!/bin/bash

echo "Ensure we are running in the proper folder"
if !([[ -d "./src/docker" ]] && [[ -f "./src/main.go" ]]); then 
    echo "This script must be run from the repository root folder"
    exit 1
fi

# Build Docker Image
docker build --file src/docker/Dockerfile --tag ghcr.io/agile-learning-institute/mentorhub-person-api:latest .
if [ $? -ne 0 ]; then
    echo "Docker build failed"
    exit 1
fi

# Start the API container with backing database
mh up person-api
