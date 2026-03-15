# =============================================================================
# Stage 1: Builder
# =============================================================================
FROM golang:1.25-bookworm AS builder

WORKDIR /app

# Install Node.js (for Tailwind CSS v4) and templ CLI dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    ca-certificates \
    && curl -fsSL https://deb.nodesource.com/setup_22.x | bash - \
    && apt-get install -y --no-install-recommends nodejs \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Install templ CLI
RUN go install github.com/a-h/templ/cmd/templ@latest

# Install Tailwind CSS v4 CLI globally via npm
RUN npm install -g @tailwindcss/cli

# ── Go dependencies (cached layer) ───────────────────────────────────────────
COPY go.mod go.sum ./
RUN go mod download

# ── Node dependencies (cached layer) ─────────────────────────────────────────
COPY package.json package-lock.json* ./
RUN npm install --ignore-scripts 2>/dev/null || true

# ── Copy full source ──────────────────────────────────────────────────────────
COPY . .

# ── Generate templ components ─────────────────────────────────────────────────
RUN templ generate

# ── Build Tailwind CSS v4 ─────────────────────────────────────────────────────
RUN tailwindcss -i ./assets/css/input.css -o ./assets/css/output.css --minify

# ── Compile Go binary ─────────────────────────────────────────────────────────
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o ./bin/main ./cmd/web/main.go


# =============================================================================
# Stage 2: Runtime
# =============================================================================
FROM debian:bookworm-slim AS runtime

WORKDIR /app

# Install FFmpeg (video thumbnails) + CA certs (HTTPS/S3 calls)
RUN apt-get update && apt-get install -y --no-install-recommends \
    ffmpeg \
    ca-certificates \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Create a non-root user to run the app
RUN useradd -r -u 1001 -g root appuser

# Copy compiled binary
COPY --from=builder /app/bin/main ./bin/main

# Copy static assets (JS, CSS, images etc.)
COPY --from=builder /app/assets ./assets

# Copy migrations (Goose runs at startup if configured)
COPY --from=builder /app/migrations ./migrations

# Ensure binary is executable
RUN chmod +x ./bin/main

# Switch to non-root user
USER appuser

EXPOSE 8080

CMD ["./bin/main"]
