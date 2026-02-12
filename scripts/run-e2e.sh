#!/bin/bash
set -e

# Ensure we are in the project root
cd "$(dirname "$0")/.."

source "$(dirname "$0")/load-env.sh"

APP_ID="agent-dev-environment"
export E2E_SERVER_URL="http://localhost:8080"
TEST_DIR="/tmp/agent-dev-environment-e2e-tests"

cleanup() {
  echo "Stopping API..."
  # This now correctly kills the API, even though output is piped
  kill $API_PID 2>/dev/null || true
  echo "Cleaning up test directory..."
  rm -rf "$TEST_DIR"
}
trap cleanup EXIT

# 1. Prepare test directory
mkdir -p "$TEST_DIR"

# 2. Start the API with prefix
# We use '> >(...)' so that $! still captures the PID of agent-api, not sed.
# We use 'sed -u' (unbuffered) so logs appear instantly.
echo "Starting API..."
export AGENT_DEV_ENVIRONMENT_LOGGING_TYPE=plain
./bin/agent-dev-environment > >(sed -u "s/^/[$APP_ID] /") 2>&1 &
API_PID=$!

# 3. Wait for the API to be ready (Loop until connection succeeds)
echo "Waiting for API to be ready..."
for i in {1..30}; do
  if curl -s http://localhost:8080/health > /dev/null; then
    echo "API is up!"
    break
  fi
  sleep 0.5
done

# 4. Run the blackbox tests
# Pass the API URL so tests know where to hit
echo "Running E2E tests..."
export API_URL="http://localhost:8080"
go test ./e2e/... -v -count=1 2>&1 | sed -u "s/^/[e2e-tests] /"
