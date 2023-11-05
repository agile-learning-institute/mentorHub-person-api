# Build stage
FROM golang:1.19 AS builder
WORKDIR /app

# Ensure git is installed
RUN apt-get update && apt-get install -y git

# Copy the entire repository into the build stage
COPY . .

# Build Go binary
RUN GOOS=linux GOARCH=amd64 go build -o institute-person-api-amd64 ./src/main.go
RUN GOOS=linux GOARCH=arm64 go build -o institute-person-api-arm64 ./src/main.go

# Get branch and patch level, then create PATCH_LEVEL file
RUN BRANCH=$(git rev-parse --abbrev-ref HEAD) && \
    PATCH=$(git rev-parse HEAD) && \
    echo $BRANCH.$PATCH > PATCH_LEVEL

# Deployment stage
FROM ubuntu
COPY --from=builder /app/institute-person-api-amd64 /institute-person-api-amd64
COPY --from=builder /app/institute-person-api-arm64 /institute-person-api-arm64
COPY --from=builder /app/PATCH_LEVEL /PATCH_LEVEL
COPY --from=builder /app/src/docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8081
CMD ["/entrypoint.sh"]
