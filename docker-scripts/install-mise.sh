#!/bin/bash
set -e

echo "Installing Mise..."
curl https://mise.jtx.dev/install.sh | sh

# Add to path for the remainder of the build process if needed
export PATH="$HOME/.local/share/mise/bin:$HOME/.local/share/mise/shims:$PATH"

echo "Mise installation complete."
