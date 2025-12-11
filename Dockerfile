# Multi-stage Dockerfile for gRPC Go application

# Stage 1: Proto generation and build
FROM golang:1.24-alpine AS builder

# Install required packages for protobuf compilation
RUN apk add --no-cache \
    protobuf \
    protobuf-dev \
    git \
    make

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Install protoc plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest

# Setup Google API proto files
RUN mkdir -p ~/.proto && \
    git clone --depth 1 --filter=blob:none --sparse https://github.com/googleapis/googleapis.git /tmp/googleapis && \
    cd /tmp/googleapis && \
    git sparse-checkout set google/api google/rpc && \
    cp -r google ~/.proto/ && \
    rm -rf /tmp/googleapis

# Copy source files
COPY . .

# Generate protobuf files
RUN make proto

# Build the application
RUN make build

# Stage 2: Final runtime image
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Change ownership to non-root user
RUN chown appuser:appgroup main

# Switch to non-root user
USER appuser

# Expose ports
EXPOSE 8080 50051

# Run the binary
CMD ["./main"]
