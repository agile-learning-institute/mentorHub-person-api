#!/bin/sh

# Determine the architecture of the current machine
ARCH=$(uname -m)

# Select the executable based on the architecture
case "$ARCH" in
    x86_64)
        exec /institute-person-api-amd64
        ;;
    aarch64)
        exec /institute-person-api-arm64
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac
