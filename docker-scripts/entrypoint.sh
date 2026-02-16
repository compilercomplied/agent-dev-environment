#!/bin/bash
set -e

if [ -z "$GITHUB_TOKEN" ]; then
    echo "Error: GITHUB_TOKEN environment variable is not set."
    exit 1
fi

if [ -z "$PULUMI_CONFIG_PASSPHRASE" ]; then
    echo "Error: PULUMI_CONFIG_PASSPHRASE environment variable is not set."
    exit 1
fi


git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"

eval "$(mise activate bash)"

exec ./agent-dev-environment
