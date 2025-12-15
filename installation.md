# Installation Guide

This guide covers the installation of Protocol Buffers (protobuf), gRPC tools, and validation tools for Go development.

## Prerequisites

- Go 1.24.0 or later
- Protocol Buffers compiler (protoc)
- Git (for downloading proto dependencies)

## Installation Steps

### 1. Install Protocol Buffers Compiler (protoc)

Install protoc on macOS using Homebrew:

```bash
brew install protobuf
```

Verify the installation:

```bash
protoc --version
```

Expected output: `libprotoc 29.3` (or later)

### 2. Install protoc-gen-go

This plugin generates Go code from `.proto` files:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

Verify the installation:

```bash
protoc-gen-go --version
```

Expected output: `protoc-gen-go v1.36.10` (or later)

### 3. Install protoc-gen-go-grpc

This plugin generates gRPC service code for Go:

```bash
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Verify the installation:

```bash
protoc-gen-go-grpc --version
```

Expected output: `protoc-gen-go-grpc 1.6.0` (or later)

### 4. Install protoc-gen-grpc-gateway

This plugin generates gRPC-Gateway reverse-proxy code:

```bash
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
```

Verify the installation:

```bash
protoc-gen-grpc-gateway --version
```

Expected output: Version information for grpc-gateway v2.27.3 (or later)

### 5. Install protoc-gen-openapiv2 (Optional)

This plugin generates OpenAPI v2 (Swagger) documentation:

```bash
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
```

### 6. Install protoc-gen-validate

This plugin generates validation code from proto validation annotations:

```bash
# Fix Go toolchain version mismatch if needed
export PATH="/usr/local/go/bin:$PATH"

# Install protoc-gen-validate
go install github.com/envoyproxy/protoc-gen-validate@v1.0.4
```

### 7. Install protoc-gen-gorm

This plugin generates GORM models from proto GORM annotations:

```bash
# Install protoc-gen-gorm
go install github.com/infobloxopen/protoc-gen-gorm@latest
```

Verify the installation:

```bash
protoc-gen-gorm --version
```

### 8. Setup GORM Proto Files

Download the gorm.proto file for protoc-gen-gorm annotations:

```bash
# Create third_party directory for GORM proto
mkdir -p third_party/github.com/infobloxopen/protoc-gen-gorm/options

# Download gorm.proto from protoc-gen-gorm repository
curl -o third_party/github.com/infobloxopen/protoc-gen-gorm/options/gorm.proto https://raw.githubusercontent.com/infobloxopen/protoc-gen-gorm/v1.1.5/proto/options/gorm.proto
```

Verify the installation:

```bash
protoc-gen-validate --version
```

**Troubleshooting Go Toolchain Mismatch:**

If you encounter errors like "compile: version does not match go tool version", this is due to Go toolchain version mismatch. Fix it by:

1. **Use system Go installation:**
   ```bash
   export PATH="/usr/local/go/bin:$PATH"
   go install github.com/envoyproxy/protoc-gen-validate@v1.0.4
   ```

2. **Verify proper Go version:**
   ```bash
   which go  # Should show /usr/local/go/bin/go
   go version  # Should match your system Go version
   ```

3. **Update your shell profile:**
   ```bash
   echo 'export PATH="/usr/local/go/bin:$PATH"' >> ~/.zshrc
   source ~/.zshrc
   ```

### 9. Setup Validation Proto Files

Download the validate.proto file for protoc-gen-validate annotations:

```bash
# Create third_party directory for validation proto
mkdir -p third_party/validate

# Download validate.proto from protoc-gen-validate repository
curl -o third_party/validate/validate.proto https://raw.githubusercontent.com/envoyproxy/protoc-gen-validate/v1.0.4/validate/validate.proto
```

### 10. Setup Google API Annotations Proto Files

To use `import "google/api/annotations.proto"` in your proto files, set up a shared proto directory in your home folder that can be used across all projects.

**Option A: Download from official googleapis repository (Recommended)**

```bash
# Clone googleapis repository (only google/api directory)
git clone --depth 1 --filter=blob:none --sparse https://github.com/googleapis/googleapis.git /tmp/googleapis
cd /tmp/googleapis
git sparse-checkout set google/api google/rpc

# Create shared proto directory and copy files
mkdir -p ~/.proto
cp -r google ~/.proto/

# Cleanup
cd -
rm -rf /tmp/googleapis
```

**Option B: Copy from grpc-gateway module cache**

```bash
# Create shared proto directory in home folder
mkdir -p ~/.proto

# Find and copy from grpc-gateway module (uses older version with bundled protos)
GRPC_GATEWAY_PATH=$(go env GOMODCACHE)/github.com/grpc-ecosystem/grpc-gateway/v2@v2.0.1
cp -r $GRPC_GATEWAY_PATH/third_party/googleapis/google ~/.proto/
```

**Note:** Newer versions of grpc-gateway (v2.20+) use Buf for dependency management and don't bundle proto files. If you need the proto files from the module cache, use an older version like v2.0.1 that includes them in `third_party/googleapis/`.

The `~/.proto/google/api/` directory now contains:
- `annotations.proto` - HTTP/gRPC annotations
- `http.proto` - HTTP configuration
- `httpbody.proto` - HTTP body definitions

When running `protoc` in any project, include the shared proto directory and third_party in your proto path:

```bash
# Ensure proper Go toolchain is used
export PATH="/usr/local/go/bin:$PATH"

# Generate protobuf code with validation and GORM models
protoc -I. -I./third_party -I$HOME/.proto \
  --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
  --validate_out="lang=go:." \
  --gorm_out=. \
  your_service.proto
```

## Verify Installation

Ensure all tools are in your PATH. The Go install command places binaries in `$GOPATH/bin` (or `$HOME/go/bin` by default).

Add to your shell profile if needed:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

## Installed Versions

Current compatible versions:

- **protoc**: libprotoc 29.3
- **protoc-gen-go**: v1.36.10
- **protoc-gen-go-grpc**: v1.6.0
- **protoc-gen-grpc-gateway**: v2.27.3
- **protoc-gen-openapiv2**: v2.27.3
- **protoc-gen-validate**: v1.0.4
- **protoc-gen-gorm**: v1.1.5

## Compatibility

All installed versions are compatible and work together seamlessly for generating:
- Protocol Buffer message definitions (Go structs)
- gRPC service implementations
- gRPC-Gateway HTTP/JSON reverse proxies
- Validation code from proto annotations
- GORM database models from proto annotations

## Next Steps

After installation, you can:
1. Define your `.proto` files with validation and GORM annotations
2. Generate Go code using `protoc` with the installed plugins
3. Implement your gRPC services with automatic validation and database models
4. Build REST APIs using gRPC-Gateway
5. Use generated GORM models with `ToORM()` and `ToPB()` conversion methods

## Proto Annotations Setup

### Validation Rules

To use validation in your proto files, add validation rules using the `validate.rules` extension:

```protobuf
syntax = "proto3";
package myservice;

import "validate/validate.proto";

message CreateUserRequest {
  string name = 1 [(validate.rules).string = {min_len: 2, max_len: 100}];
  string email = 2 [(validate.rules).string = {pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"}];
}
```

### GORM Model Annotations

To generate GORM models from your proto files, add GORM annotations:

```protobuf
syntax = "proto3";
package myservice;

import "github.com/infobloxopen/protoc-gen-gorm/options/gorm.proto";

message UserEntity {
  option (gorm.opts) = {
    ormable: true,
    table: "users"
  };
  
  uint32 id = 1 [(gorm.field).tag = {primary_key: true, auto_increment: true}];
  string name = 2 [(gorm.field).tag = {not_null: true, size: 100}];
  string email = 3 [(gorm.field).tag = {unique_index: "email_idx", size: 255}];
}
```

### Combined Usage

You can use both validation and GORM annotations together:

```protobuf
syntax = "proto3";
package myservice;

import "validate/validate.proto";
import "github.com/infobloxopen/protoc-gen-gorm/options/gorm.proto";

message UserEntity {
  option (gorm.opts) = {ormable: true, table: "users"};
  
  uint32 id = 1 [(gorm.field).tag = {primary_key: true, auto_increment: true}];
  string name = 2 [
    (validate.rules).string = {min_len: 2, max_len: 100},
    (gorm.field).tag = {not_null: true, size: 100}
  ];
  string email = 3 [
    (validate.rules).string = {pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"},
    (gorm.field).tag = {unique_index: "email_idx", size: 255}
  ];
}
```

This approach provides both automatic validation and database model generation from a single proto definition.
