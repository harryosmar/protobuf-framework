# Protobuf Go gRPC Server

A gRPC server with HTTP gateway demonstrating Protocol Buffers and gRPC in Go with multiple services.

## Features

- **gRPC Services**: 
  - `HelloService` with `GetHello` RPC method
  - `UserService` with `CreateUser` and `GetUser` RPC methods
- **HTTP Gateway**: REST API endpoints using grpc-gateway
- **Protocol Buffers**: Message definitions for Hello and User services
- **Multi-stage Docker build**: Optimized containerization

## Project Structure

```
.
├── proto/              # Protocol buffer definitions
│   ├── hello.proto
│   └── user.proto
├── gen/                # Generated code from proto files (gitignored)
│   ├── hello/
│   │   ├── hello.pb.go
│   │   ├── hello_grpc.pb.go
│   │   └── hello.pb.gw.go
│   └── user/
│       ├── user.pb.go
│       ├── user_grpc.pb.go
│       └── user.pb.gw.go
├── service/            # Service implementations
│   ├── HelloService.go
│   ├── UserService.go
│   └── service.go
├── scripts/            # Build scripts
│   └── generate.sh
├── main.go             # Main application
├── Makefile            # Build automation
├── Dockerfile          # Multi-stage Docker build
└── installation.md     # Installation guide for tools
```

## Prerequisites

See [installation.md](installation.md) for detailed setup instructions.

Required tools:
- Go 1.24.0+
- protoc (Protocol Buffers compiler)
- protoc-gen-go
- protoc-gen-go-grpc
- protoc-gen-grpc-gateway

## Running the Server

### Local Development

```bash
# Generate protobuf files
make proto

# Run the server
make run
# or
go run main.go
```

### Using Docker

```bash
# Build the Docker image
docker build -t protobuf-go-app .

# Run the container
docker run -p 8080:8080 -p 50051:50051 protobuf-go-app
```

The server will start:
- **gRPC server** on port `50051`
- **HTTP gateway** on port `8080`

### Test with gRPC

Using grpcurl:

```bash
grpcurl -plaintext -d '{"name": "World"}' localhost:50051 hello.HelloService/GetHello
```

### Test with HTTP/REST

**HelloService:**
```bash
curl http://localhost:8080/v1/hello/World
```

**UserService:**
```bash
# Create a user
curl -X POST http://localhost:8080/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}'

# Get a user
curl http://localhost:8080/v1/users/1
```

## Architecture

### Request Flow

The server supports both gRPC and HTTP protocols with a unified middleware architecture:

```
┌─────────────┐    ┌──────────────┐    ┌─────────────┐    ┌─────────────┐
│ gRPC Client │───▶│ gRPC Server  │───▶│ Interceptor │───▶│   Service   │
└─────────────┘    │   :50051     │    │ (RequestID) │    │ (Hello/User)│
                   └──────────────┘    └─────────────┘    └─────────────┘

┌─────────────┐    ┌──────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│ HTTP Client │───▶│ HTTP Gateway │───▶│ gRPC Server │───▶│ Interceptor │───▶│   Service   │
└─────────────┘    │   :8080      │    │   :50051    │    │ (RequestID) │    │ (Hello/User)│
                   └──────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
```

### Middleware Strategy

- **Single gRPC Interceptor**: Handles request ID generation for both direct gRPC and HTTP gateway requests
- **Automatic Propagation**: gRPC-Gateway forwards gRPC metadata to HTTP response headers
- **Consistent Behavior**: Both protocols receive `X-Request-ID` headers with UUID values
- **Request Logging**: All requests logged with unique identifiers for tracing

### Benefits

- **DRY Principle**: One middleware handles both protocols
- **Maintainability**: Single source of truth for request ID logic
- **Performance**: No duplicate processing for gateway requests
- **Observability**: Consistent request tracing across all client types

## Build Commands

### Available Make Targets

```bash
make proto   # Generate protobuf files from .proto sources
make build   # Build static binary for production
make clean   # Remove generated files
make run     # Run development server
```

### Manual Proto Generation

If you prefer to generate manually:

```bash
make proto
```

## API Documentation

### HelloService

#### GetHello

**gRPC Method**: `hello.HelloService/GetHello`

**HTTP Endpoint**: `GET /v1/hello/{name}`

**Request**:
```protobuf
message HelloRequest {
  string name = 1;
}
```

**Response**:
```protobuf
message HelloResponse {
  string message = 1;
}
```

**Example**:
- Request: `{"name": "Alice"}`
- Response: `{"message": "Hello, Alice!"}`

### UserService

#### CreateUser

**gRPC Method**: `user.UserService/CreateUser`

**HTTP Endpoint**: `POST /v1/users`

**Request**:
```protobuf
message CreateUserRequest {
  UserDTO user = 1;
}

message UserDTO {
  string name = 1;
  string email = 2;
}
```

**Response**:
```protobuf
message CreateUserResponse {
  UserEntity user = 1;
}

message UserEntity {
  int64 id = 1;
  string name = 2;
  string email = 3;
  string created_at = 4;
  string updated_at = 5;
}
```

#### GetUser

**gRPC Method**: `user.UserService/GetUser`

**HTTP Endpoint**: `GET /v1/users/{id}`

**Request**:
```protobuf
message GetUserRequest {
  int64 id = 1;
}
```

**Response**:
```protobuf
message GetUserResponse {
  UserEntity user = 1;
}
```
