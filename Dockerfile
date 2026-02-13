# Base stage with common dev tools
FROM ubuntu:24.04 AS base
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y \
    curl \
    git \
    ca-certificates \
    wget \
    ripgrep \
    && rm -rf /var/lib/apt/lists/*

# Builder stage
FROM base AS builder
# Install build-specific tools
RUN apt-get update && apt-get install -y build-essential && rm -rf /var/lib/apt/lists/*

# Copy and run Go installation script
COPY docker-scripts/install-go.sh /tmp/install-go.sh
RUN chmod +x /tmp/install-go.sh && /tmp/install-go.sh && rm /tmp/install-go.sh
ENV PATH="/usr/local/go/bin:$PATH"

WORKDIR /app
COPY go.mod ./
COPY go.sum* ./
RUN go mod download

COPY . .
RUN go build -o agent-dev-environment ./src

# Runner stage (The "Fat" Dev Environment)
FROM base AS runner
WORKDIR /app

# Install mise in the runner stage
RUN curl https://mise.jtx.dev/install.sh | sh
ENV PATH="/root/.local/share/mise/bin:/root/.local/share/mise/shims:$PATH"

# Copy the compiled binary from builder
COPY --from=builder /app/agent-dev-environment .

EXPOSE 8080

CMD ["./agent-dev-environment"]
