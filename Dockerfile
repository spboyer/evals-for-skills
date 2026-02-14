# Dockerfile for waza - AI agent skill evaluation CLI
# This provides a containerized environment for running waza in CI/CD pipelines

FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN go build -v -o waza ./cmd/waza

# Runtime stage - minimal image
FROM alpine:latest

RUN apk add --no-cache ca-certificates git

WORKDIR /workspace

# Copy the binary from builder
COPY --from=builder /build/waza /usr/local/bin/waza

# Verify installation
RUN waza --version

# Default command
ENTRYPOINT ["waza"]
CMD ["--help"]
