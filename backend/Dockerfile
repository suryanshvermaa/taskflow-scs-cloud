# syntax=docker/dockerfile:1

# ── Stage 1: Build ────────────────────────────────────────────────────────────
# Use Debian-based Go image for the builder so CGO tools (gcc) are always
# available without relying on Alpine's apk package index.
FROM golang:1.27-bookworm AS builder

WORKDIR /src

# Cache dependencies first.
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build a static binary.
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build \
        -ldflags="-s -w" \
        -o /bin/taskflow-api \
        ./cmd/server

# ── Stage 2: Runtime ──────────────────────────────────────────────────────────
# Keep the final image small with Debian slim.
FROM debian:bookworm-slim

# ca-certificates for outbound HTTPS if ever needed.
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Non-root user for security.
RUN groupadd -r app && useradd -r -g app app

WORKDIR /app

# Persistent data directory (mount a volume here in production).
RUN mkdir -p /app/data && chown app:app /app/data

COPY --from=builder /bin/taskflow-api /app/taskflow-api
RUN chown app:app /app/taskflow-api

USER app

EXPOSE 8080

ENV PORT=8080
ENV DB_PATH=/app/data/taskflow.db
ENV CORS_ORIGIN=http://localhost:5173

ENTRYPOINT ["/app/taskflow-api"]
