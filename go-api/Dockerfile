# Multi-stage build for optimized Go API image

# Stage 1: Build
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod file
COPY go.mod ./

# Download dependencies (generates go.sum)
RUN go mod download

# Copy source code
COPY . .

# Tidy dependencies and build (go build will verify checksums in go.sum)
RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Stage 2: Runtime
FROM alpine:latest

# Install CA certificates and wget for health checks
RUN apk --no-cache add ca-certificates tzdata wget

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/main .

# Create non-root user (using 1001 to avoid conflict with default GID 1000)
RUN addgroup -g 1001 appuser && \
    adduser -D -u 1001 -G appuser appuser && \
    chown -R appuser:appuser /app

USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health 2>/dev/null || exit 1

# Run the application
CMD ["./main"]
