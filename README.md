# Agent Dev Environment

The "hands" of a coding agent. This service provides a sandboxed Ubuntu environment with a Go HTTP API that gives an AI agent the ability to interact with a filesystem and development tools. The agent's "brains" (the `agent-hub` component) connects to this environment to read, write, search, and manage files — as well as run commands like `git`, build tools, and linters — all abstracted through [mise](https://mise.jdx.dev/).

## Architecture

```
┌──────────────┐         HTTP          ┌─────────────────────────────────┐
│   agent-hub  │ ───────────────────▶  │   agent-dev-environment (pod)   │
│   (brains)   │                       │                                 │
└──────────────┘                       │  ┌───────────────────────────┐  │
                                       │  │  Go HTTP API (:8080)      │  │
                                       │  │  • filesystem operations  │  │
                                       │  │  • command execution      │  │
                                       │  └───────────────────────────┘  │
                                       │                                 │
                                       │  Ubuntu 24.04 image with:      │
                                       │  • git, curl, wget             │
                                       │  • build-essential             │
                                       │  • mise (task runner)          │
                                       │  • Go 1.25                     │
                                       └─────────────────────────────────┘
```

The container ships as a fat Ubuntu-based image that includes the compiled Go API server alongside a full development toolchain. When the pod starts, the API server launches on port 8080 and the `agent-hub` begins orchestrating work against it.

## API

All endpoints accept and return JSON.

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Health check — returns `"OK"` |
| `POST` | `/api/v1/filesystem/read` | Read file content |
| `POST` | `/api/v1/filesystem/create_file` | Create a new file (with auto-directory creation) |

### Examples

**Read a file:**
```json
// POST /api/v1/filesystem/read
{ "path": "/workspace/main.go" }

// 200 OK
{ "content": "package main\n..." }
```

**Create a file:**
```json
// POST /api/v1/filesystem/create_file
{ "path": "/workspace/hello.txt", "content": "Hello, world!" }

// 200 OK
{}
```

Errors follow a consistent format:
```json
{ "status": 404, "message": "file not found" }
```

## Project Structure

```
├── src/                        # Go HTTP API
│   ├── main.go                 # Server entry point and route registration
│   ├── healthchecks.go         # Health check handler
│   ├── api/v1/                 # Request/response models
│   ├── features/               # Route handlers (business logic)
│   │   └── filesystem/         # Filesystem operations (read, create_file)
│   ├── internal/middleware/     # Panic recovery middleware
│   └── library/                # Shared utilities
│       ├── api/                # Generic typed handler wrapper
│       ├── config/             # Environment variable management
│       └── logger/             # Plain or structured (JSON) logging
├── e2e/                        # End-to-end tests
├── iac/                        # Pulumi infrastructure (Kubernetes)
├── scripts/                    # Helper scripts
├── docker-scripts/             # Scripts used during Docker build
├── Dockerfile                  # Production image (Ubuntu 24.04)
├── Dockerfile.test             # Lightweight image for running e2e tests
├── docker-compose.e2e.yaml     # Compose file for local e2e testing
└── mise.toml                   # Task runner and tool version management
```

## Mise

[mise](https://mise.jdx.dev/) is used to manage tool versions and abstract common tasks. It is installed in the Docker image and available at runtime.

**Managed tools:**

| Tool | Version |
|------|---------|
| Go | 1.25.0 |
| Node | 24 |
| Pulumi | latest |
| gopls | 0.20.0 |

**Tasks:**

| Task | Command | Description |
|------|---------|-------------|
| `mise run build` | `go build -o bin/agent-orchestrator ./src` | Compile the Go API |
| `mise run run` | `go run ./src` | Run the API locally |
| `mise run test:e2e` | `go test ./e2e/...` | Run e2e tests against a running server |
| `mise run test:e2e-docker` | `docker compose ...` | Run e2e tests via Docker Compose |

## Development

### Prerequisites

- [mise](https://mise.jdx.dev/) installed, or Go 1.25+ and Docker
- Docker and Docker Compose (for e2e tests)

### Run locally

```bash
# Install tools via mise
mise install

# Start the server
AGENT_DEV_ENVIRONMENT_LOGGING_TYPE=plain mise run run

# In another terminal, run the e2e tests
E2E_SERVER_URL=http://localhost:8080 mise run test:e2e
```

### Run e2e tests with Docker

```bash
mise run test:e2e-docker
# or directly:
./scripts/run-e2e.sh
```

This builds the app image, starts the container with a health check, runs the test suite, and tears everything down.

## Deployment

### Docker Image

The production image is a multi-stage build (`Dockerfile`) based on Ubuntu 24.04. The builder stage compiles the Go binary, and the runner stage packages it alongside the full development toolchain (git, build-essential, mise, Go). The image is pushed to **GitHub Container Registry** (`ghcr.io`).

### CI/CD

A single GitHub Actions workflow (`.github/workflows/cicd.yaml`) handles everything:

1. **On every push and PR to `master`:** builds the Docker image and runs the e2e test suite.
2. **On push to `master` only:** builds and pushes the image to `ghcr.io`, then deploys to Kubernetes via Pulumi.

### Infrastructure

Infrastructure is managed with [Pulumi](https://www.pulumi.com/) (TypeScript) in the `iac/` directory. It provisions:

- A **Kubernetes ConfigMap** (`agent-dev-env-config`) for plain configuration values
- A **Kubernetes Secret** (`agent-dev-env-secret`) for sensitive configuration values

Configuration keys prefixed with `AGENT_DEV_ENVIRONMENT_` are automatically discovered from the Pulumi stack config and routed to the appropriate resource.

**Stack configurations:**

| Stack | Logging | File |
|-------|---------|------|
| `local` | `plain` | `iac/Pulumi.local.yaml` |
| `prod` | `structured` | `iac/Pulumi.prod.yaml` |

## Configuration

The application is configured via environment variables prefixed with `AGENT_DEV_ENVIRONMENT_`:

| Variable | Required | Values | Description |
|----------|----------|--------|-------------|
| `AGENT_DEV_ENVIRONMENT_LOGGING_TYPE` | Yes | `plain`, `structured` | Log output format |
