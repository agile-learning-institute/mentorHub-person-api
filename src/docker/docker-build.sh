#!/bin/bash

# Build Go binary
# Build Docker image
docker build --file src/docker/Dockerfile --tag ghcr.io/agile-learning-institute/institute-person-api:latest .
if [ $? -ne 0 ]; then
    echo "Docker build failed"
    exit 1
fi

# Push Docker image
if [ $1 = '--push' ]; then
    docker push ghcr.io/agile-learning-institute/institute-person-api:latest
    if [ $? -ne 0 ]; then
    echo "Docker build failed"
    exit 1
        echo "Docker push failed"
        exit 1
    fi
    echo "image pushed"
fi

# Run the Database and API containers
if [ $1 = '--run' ]; then
    curl https://raw.githubusercontent.com/agile-learning-institute/institute/main/docker-compose/run-local-person-api.sh | /bin/bash
fi
