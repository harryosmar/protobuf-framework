#!/bin/bash

# Exit on error
set -e

# Create output directory if it doesn't exist
mkdir -p gen

# Generate code from all proto files
protoc -I./proto -I$HOME/.proto \
  --go_out=./gen --go_opt=paths=source_relative \
  --go-grpc_out=./gen --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=./gen --grpc-gateway_opt=paths=source_relative \
  proto/*.proto

echo "✓ Proto files generated successfully"
