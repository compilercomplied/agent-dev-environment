# --- 1. BASE TOOLSET ---
FROM ubuntu:24.04 AS base
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y \
    curl git ca-certificates wget ripgrep \
    && rm -rf /var/lib/apt/lists/*

# --- 2. BUILDER (App & SDKs) ---
FROM base AS builder
RUN apt-get update && apt-get install -y build-essential && rm -rf /var/lib/apt/lists/*
COPY docker-scripts/install-go.sh /tmp/install-go.sh
RUN chmod +x /tmp/install-go.sh && /tmp/install-go.sh && rm /tmp/install-go.sh
ENV PATH="/usr/local/go/bin:$PATH"

WORKDIR /app
COPY go.mod ./
COPY go.sum* ./
RUN go mod download
COPY . .
RUN go build -o agent-dev-environment ./src

# --- 3. OPENAPI SPEC (Export Target) ---
FROM builder AS spec-gen
RUN go run github.com/swaggo/swag/cmd/swag@latest init -g src/main.go -o ./docs --parseDependency --parseInternal

FROM scratch AS openapi-spec
COPY --from=spec-gen /app/docs/swagger.json /swagger.json

# --- 4. RUNNER (Production Env) ---
FROM base AS runner

# Create a non-root user
RUN useradd -m -s /bin/bash agent
USER agent
WORKDIR /home/agent/app

# Install mise as the non-root user
COPY --chown=agent:agent docker-scripts/install-mise.sh /tmp/install-mise.sh
RUN chmod +x /tmp/install-mise.sh && /tmp/install-mise.sh && rm /tmp/install-mise.sh

# Update PATH for the agent user
ENV PATH="/home/agent/.local/share/mise/bin:/home/agent/.local/share/mise/shims:$PATH"

COPY --chown=agent:agent --from=builder /app/agent-dev-environment .
COPY --chown=agent:agent docker-scripts/entrypoint.sh /home/agent/app/entrypoint.sh
RUN chmod +x /home/agent/app/entrypoint.sh

EXPOSE 8080
ENTRYPOINT ["/home/agent/app/entrypoint.sh"]
