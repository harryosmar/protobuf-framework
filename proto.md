# Protocol Buffer Code Generation Guide

This guide explains how to generate Go code from `.proto` files using `protoc`.

## Generated Files

When you run `protoc` on a `.proto` file, it generates three types of Go files:

1. **`*.pb.go`** - Protocol Buffer message definitions (structs)
2. **`*_grpc.pb.go`** - gRPC service client and server interfaces
3. **`*.pb.gw.go`** - gRPC-Gateway HTTP reverse-proxy handlers

## Prerequisites

Ensure you have installed all required tools (see [installation.md](installation.md)):
- `protoc` - Protocol Buffers compiler
- `protoc-gen-go` - Go code generator
- `protoc-gen-go-grpc` - gRPC service generator
- `protoc-gen-grpc-gateway` - gRPC-Gateway generator

## Basic Command

```bash
protoc -I./proto -I$HOME/.proto \
  --go_out=./gen --go_opt=paths=source_relative \
  --go-grpc_out=./gen --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=./gen --grpc-gateway_opt=paths=source_relative \
  proto/hello.proto
```

## Command Breakdown

### Include Paths (`-I` flags)

```bash
-I./proto           # Include the proto directory
-I$HOME/.proto      # Include shared proto files (google/api/annotations.proto)
```

- `-I` specifies directories where `protoc` looks for imported `.proto` files
- Multiple `-I` flags can be used to specify multiple include paths

### Go Code Generation (`--go_out`)

```bash
--go_out=./gen                    # Output directory for .pb.go files
--go_opt=paths=source_relative    # Keep same directory structure as source
```

Generates: `hello.pb.go`
- Contains Go structs for `HelloRequest` and `HelloResponse`
- Implements protobuf serialization/deserialization

### gRPC Service Generation (`--go-grpc_out`)

```bash
--go-grpc_out=./gen                    # Output directory for _grpc.pb.go files
--go-grpc_opt=paths=source_relative    # Keep same directory structure as source
```

Generates: `hello_grpc.pb.go`
- Contains `HelloServiceServer` interface
- Contains `HelloServiceClient` interface
- Implements gRPC service registration

### gRPC-Gateway Generation (`--grpc-gateway_out`)

```bash
--grpc-gateway_out=./gen                    # Output directory for .pb.gw.go files
--grpc-gateway_opt=paths=source_relative    # Keep same directory structure as source
```

Generates: `hello.pb.gw.go`
- Contains HTTP handler registration
- Converts HTTP/JSON requests to gRPC calls
- Enables REST API endpoints

## Generate Script

```bash
# Generate code from all proto files
protoc -I./proto -I$HOME/.proto \
  --go_out=./gen --go_opt=paths=source_relative \
  --go-grpc_out=./gen --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=./gen --grpc-gateway_opt=paths=source_relative \
  proto/*.proto
```

## Common Options

### Path Options

- `paths=source_relative` - Output files in same directory structure as source
- `paths=import` - Output files based on Go import path

### Module Path

If using `paths=import`, specify the module path:

```bash
--go_opt=module=github.com/harryosmar/protobuf-go
```

## Troubleshooting

### Error: "google/api/annotations.proto: File not found"

**Solution**: Ensure you have set up the Google API proto files:

```bash
mkdir -p ~/.proto
# Follow instructions in installation.md section 6
```

### Error: "protoc-gen-go: program not found"

**Solution**: Ensure the plugin is installed and in your PATH:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Error: "No such file or directory: ./gen"

**Solution**: Create the output directory first:

```bash
mkdir -p gen
```

## Best Practices

1. **Version control**: Add generated files to `.gitignore` and regenerate on build
2. **Automation**: Use Makefile or shell scripts for consistent generation
3. **CI/CD**: Include proto generation in your build pipeline
4. **Documentation**: Keep proto files well-documented with comments

## Example Output

After running the command, you should see:

```
gen/
├── hello.pb.go          # Message definitions
├── hello_grpc.pb.go     # gRPC service
└── hello.pb.gw.go       # HTTP gateway
```

## Next Steps

After generating the code:

1. Run `go mod tidy` to download dependencies
2. Implement the service interface in your server code
3. Start the server with `go run cmd/server/main.go`
