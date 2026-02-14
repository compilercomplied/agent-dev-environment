# Agent Dev Environment

The "hands" of a coding agent. This service provides a sandboxed Ubuntu environment with a Go HTTP API that gives an AI agent the ability to interact with a filesystem and development tools. The agent's "brains" (the `agent-hub` component) connects to this environment to read, write, search, and manage files — as well as run commands like `git`, build tools, and linters — all abstracted through [mise](https://mise.jdx.dev/).

## Architecture

```
┌──────────────┐         HTTP          ┌─────────────────────────────────┐
│   agent-hub  │ ───────────────────▶  │   agent-dev-environment (pod)   │
│              │                       │                                 │
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

## Restricted Shell Execution

The service provides a `POST /api/v1/shell/run` endpoint to execute a controlled set of commands. 

**Allowed Commands:** `ls`, `rg`, `git`, `curl`.

**Security Restrictions:**
- `curl` is restricted to `localhost` targets only (e.g., `http://localhost:8080/health`).
- Only whitelisted commands can be executed.
- All other commands are rejected with a `400 Bad Request`.

## Mise

[mise](https://mise.jdx.dev/) is used to manage tool versions and abstract common tasks. It is installed in the Docker image and available at runtime.

## Development

### Prerequisites

- [mise](https://mise.jdx.dev/) installed.

### Run locally

```bash
# Install tools via mise
mise install

# Run E2E
mise run test:e2e
```

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
