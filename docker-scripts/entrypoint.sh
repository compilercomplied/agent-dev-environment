#!/bin/bash
set -e

if [ -z "$PULUMI_CONFIG_PASSPHRASE" ]; then
    echo "Error: PULUMI_CONFIG_PASSPHRASE environment variable is not set."
    exit 1
fi



# We can use variables that start with GITHUB_ in the pipeline.
# While this variable should be absolutely mandatory, we can't enforce it if we
# want E2E tests to work in the pipeline and gatekeep code quality.
if [ "$GITHUB_TOKEN" ]; then
	git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"
else
	echo "GITHUB_TOKEN is not set, skipping git configuration"
fi

eval "$(mise activate bash)"

exec ./agent-dev-environment
