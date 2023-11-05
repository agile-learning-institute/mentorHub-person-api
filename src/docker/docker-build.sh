#!/bin/bash

# Build Docker Image
docker build --file src/docker/Dockerfile --tag ghcr.io/agile-learning-institute/institute-person-api:latest . $1
if [ $? -ne 0 ]; then
    echo "Docker build failed"
    exit 1
fi

curl https://raw.githubusercontent.com/agile-learning-institute/institute/main/docker-compose/run-local-person-api.sh | /bin/bash
