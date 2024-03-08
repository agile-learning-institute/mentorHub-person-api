# Base build stage
FROM golang:1.19 AS base
WORKDIR /src
COPY . .

# Build argument to pass the architecture
ARG TARGETARCH

# Complie for the target architecture (limited support)
RUN GOOS=linux GOARCH=$TARGETARCH go build -o /src/entrypoint ./src/main.go

# Get branch and patch level, then create PATCH_LEVEL file
RUN BRANCH=$(git rev-parse --abbrev-ref HEAD) && \
    DATE=$(date "+%Y-%m-%d:%H:%M:%S") && \
    echo $DATE.$PATCH > /src/PATCH_LEVEL

# Final stage
FROM ubuntu
ARG TARGETARCH

# Copy the binary and PATCH_LEVEL from the build stage
COPY --from=base /src/entrypoint /entrypoint
COPY --from=base /src/PATCH_LEVEL /PATCH_LEVEL

EXPOSE 8082
CMD ["/entrypoint"]
