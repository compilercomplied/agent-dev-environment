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
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    swag init -g src/main.go -o ./docs --parseDependency --parseInternal

FROM scratch AS openapi-spec
COPY --from=spec-gen /app/docs/swagger.json /swagger.json

# --- 4. RUNNER (Production Env) ---
FROM base AS runner
WORKDIR /app
COPY docker-scripts/install-mise.sh /tmp/install-mise.sh
RUN chmod +x /tmp/install-mise.sh && /tmp/install-mise.sh && rm /tmp/install-mise.sh
ENV PATH="/root/.local/share/mise/bin:/root/.local/share/mise/shims:$PATH"

COPY --from=builder /app/agent-dev-environment .
COPY docker-scripts/entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

EXPOSE 8080
ENTRYPOINT ["/app/entrypoint.sh"]
