# Build stage
FROM golang:1.24.1 AS builder


# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with security flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /go/bin/api-gateway \
    ./cmd/main.go

# Final stage
FROM alpine:3.19

# Add non root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Install necessary runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Set working directory
WORKDIR /app

# Create config directory
RUN mkdir -p /app/config

# Copy binary from builder
COPY --from=builder /go/bin/api-gateway .

# Copy config file - make sure the path is correct
COPY --from=builder /app/config/config.yaml /app/config/

# Set ownership
RUN chown -R appuser:appgroup /app

# Use non root user
USER appuser

# Expose port
EXPOSE 8080

# Set environment variables with defaults
ENV HTTP_PORT=8080 \
    GRPC_AUTH_SERVER=auth-service:50051 \
    GRPC_DOC_SERVER=doc-service:50052 \
    CORS_ALLOWED_ORIGINS="*" \
    CORS_ALLOWED_METHODS="GET,POST,PUT,DELETE" \
    CORS_ALLOWED_HEADERS="Content-Type,Authorization"

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:${HTTP_PORT}/health || exit 1

# Command to run the application
CMD ["./api-gateway"]
