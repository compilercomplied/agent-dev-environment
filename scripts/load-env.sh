#!/bin/bash

# This script loads environment variables from Pulumi's local stack
# into the current shell session.
# Usage: source scripts/load-env.sh

if ! command -v pulumi &> /dev/null; then
    echo "Error: pulumi CLI is not installed."
    return 1
fi

if ! command -v jq &> /dev/null; then
    echo "Error: jq is required but not installed."
    return 1
fi

if [ -z "$PULUMI_CONFIG_PASSPHRASE" ]; then
    echo "Error: PULUMI_CONFIG_PASSPHRASE is not set."
    return 1
fi

# Ensure Pulumi is logged in without overriding existing sessions
if ! pulumi whoami &> /dev/null; then
    if [ -n "$PULUMI_ACCESS_TOKEN" ]; then
        echo "Logging in to Pulumi Cloud..."
        pulumi login
    else
        echo "No Pulumi session found. Falling back to local backend..."
        pulumi login --local
    fi
fi

# Ensure the local stack exists
pulumi stack select --create local -C iac

echo "Loading local configuration from Pulumi..."

# Extract configuration from Pulumi and export it
# Handles the new nested env_vars structure and legacy top-level keys
set -o pipefail
if ! pulumi config --stack local -C iac --show-secrets --json | \
jq -r --arg q "'" '
    to_entries | .[] | 
    if .key | endswith(":env_vars") then
        (.value.objectValue // (.value.value | fromjson)) | to_entries | .[] | .key + "=" + $q + (.value.default | tostring) + $q
    else
        (.key | split(":") | last) + "=" + $q + 
        (
            if .value | type == "object" 
            then .value.value 
            else .value 
            end | tostring
        ) + $q
    end
' > .env; then
    echo "Error: Failed to extract configuration from Pulumi."
    return 1
fi
set +o pipefail

# Load the .env file into the current shell
set -a
source .env
set +a

echo "Environment loaded successfully from local stack"
