#!/bin/bash
set -e

echo "Installing Mise..."
curl https://mise.jtx.dev/install.sh | sh

# Add to path for the remainder of the build process if needed
export PATH="/root/.local/share/mise/bin:/root/.local/share/mise/shims:$PATH"

echo "Mise installation complete."
