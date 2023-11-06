# Base build stage
FROM golang:1.19 AS base
WORKDIR /src
COPY . .

# Build argument to pass the architecture
ARG TARGETARCH

# Map Docker's TARGETARCH to Go's GOARCH and set GOOS
RUN if [ "${TARGETARCH}" = "amd64" ]; then \
      GOARCH="amd64"; \
    elif [ "${TARGETARCH}" = "arm64" ]; then \
      GOARCH="arm64"; \
    else \
      echo "Unsupported architecture: ${TARGETARCH}" >&2; \
      exit 1; \
    fi && \
    GOOS=linux GOARCH=$GOARCH go build -o /src/institute-person-api ./src/main.go

# Get branch and patch level, then create PATCH_LEVEL file
RUN BRANCH=$(git rev-parse --abbrev-ref HEAD) && \
    PATCH=$(git rev-parse HEAD) && \
    echo "$BRANCH.$PATCH" > /src/PATCH_LEVEL

# Final stage
FROM ubuntu
ARG TARGETARCH

# Copy the binary and PATCH_LEVEL from the build stage
COPY --from=base /src/institute-person-api /institute-person-api
COPY --from=base /src/PATCH_LEVEL /PATCH_LEVEL

EXPOSE 8081
CMD ["/institute-person-api"]
