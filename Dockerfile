# ─────────────────────────────────────────────────────────────────────────────
# Stage 1 — Builder
# Uses the full Go toolchain to compile a static binary.
# This stage is discarded after build; none of its ~800MB ends up in the image.
# ─────────────────────────────────────────────────────────────────────────────
FROM golang:1.24-alpine AS builder

# Install git (needed by go mod download for some VCS-based modules)
RUN apk add --no-cache git

WORKDIR /app

# Copy dependency manifests first so Docker can cache this layer.
# The expensive "go mod download" only re-runs when go.mod / go.sum change.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build a fully-static binary (CGO_ENABLED=0 → no libc dependency)
# -ldflags="-s -w" strips debug info, shrinking the binary ~30%
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o careercompass ./cmd/api/main.go

# ─────────────────────────────────────────────────────────────────────────────
# Stage 2 — Runtime
# Minimal Alpine image (~5MB). Only the compiled binary + migration files land here.
# Final image size: ~15–20MB total.
# ─────────────────────────────────────────────────────────────────────────────
FROM alpine:3.20

# ca-certificates → needed for HTTPS calls (Gemini API)
# tzdata         → correct timezone handling for timestamps
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy only what the app needs at runtime
COPY --from=builder /app/careercompass .
COPY --from=builder /app/migrations ./migrations

# The app reads all config from environment variables (no .env file in prod)
# Expose is not needed for Railway as it uses the $PORT variable dynamically

# Run as non-root for security (principle of least privilege)
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

ENTRYPOINT ["./careercompass"]
