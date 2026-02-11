#!/bin/bash
set -e

GO_VERSION="1.25.0"
ARCH=$(uname -m)

if [ "$ARCH" = "x86_64" ]; then
    GO_ARCH="amd64"
elif [ "$ARCH" = "aarch64" ]; then
    GO_ARCH="arm64"
else
    echo "Unsupported architecture: $ARCH"
    exit 1
fi

echo "Installing Go ${GO_VERSION} for ${GO_ARCH}..."
curl -L "https://go.dev/dl/go${GO_VERSION}.linux-${GO_ARCH}.tar.gz" | tar -C /usr/local -xzf -
