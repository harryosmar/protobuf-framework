# Protobuf Go gRPC Server

A production-ready gRPC server with HTTP gateway demonstrating Protocol Buffers, gRPC, and validation in Go.

## Features

- **gRPC Services**: 
  - `HelloService` with `GetHello` RPC method
  - `UserService` with `CreateUser` and `GetUser` RPC methods
- **HTTP Gateway**: REST API endpoints using grpc-gateway
- **Protocol Buffers**: Message definitions with validation rules
- **Validation**: protoc-gen-validate for automatic validation from proto annotations
- **GORM Models**: protoc-gen-gorm for automatic database model generation
- **Repository Pattern**: Clean architecture with dependency injection
- **Production Features**: Metrics, logging, rate limiting, graceful shutdown
- **Multi-stage Docker build**: Optimized containerization

## Project Structure

```
.
├── proto/              # Protocol buffer definitions with validation
│   ├── hello.proto
│   └── user.proto
├── third_party/        # Third-party proto files
│   ├── validate/
│   │   └── validate.proto
│   └── github.com/infobloxopen/protoc-gen-gorm/options/
│       └── gorm.proto
├── gen/                # Generated code from proto files (gitignored)
│   ├── hello/
│   │   ├── hello.pb.go
│   │   ├── hello_grpc.pb.go
│   │   ├── hello.pb.gw.go
│   │   └── hello.pb.validate.go
│   └── user/
│       ├── user.pb.go
│       ├── user_grpc.pb.go
│       ├── user.pb.gw.go
│       ├── user.pb.validate.go
│       └── user.pb.gorm.go
├── service/            # Service implementations
│   ├── HelloService.go
│   └── UserService.go
├── repository/         # Data access layer
│   ├── user_repository.go
│   └── errors.go
├── middleware/         # gRPC middleware
│   ├── logging.go
│   ├── metrics.go
│   ├── ratelimit.go
│   └── requestid.go
├── database/           # Database configuration
│   └── database.go
├── config/             # Configuration management
│   └── config.go
├── logger/             # Logging utilities
│   └── logger.go
├── handlers/           # HTTP handlers
│   ├── health.go
│   └── swagger.go
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
- protoc-gen-openapiv2 (for Swagger generation)
- protoc-gen-validate (for validation generation)
- protoc-gen-gorm (for GORM model generation)

### Install Dependencies

```bash
# Fix Go toolchain version mismatch if needed
export PATH="/usr/local/go/bin:$PATH"

# Install Protocol Buffer generators
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

# Install protoc-gen-validate for validation (specific version for compatibility)
go install github.com/envoyproxy/protoc-gen-validate@v1.0.4

# Install protoc-gen-gorm for GORM model generation
go install github.com/infobloxopen/protoc-gen-gorm@latest

# Install project dependencies
go mod download
```

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

### Additional Endpoints

**Health Check:**
```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "service_name": "protobuf-go-server",
  "version": "v1.0.0",
  "status": "healthy"
}
```

**Swagger Documentation:**
- **Swagger UI**: `http://localhost:8080/docs`
- **Swagger JSON**: `http://localhost:8080/docs/swagger.json`

**Prometheus Metrics:**
- **Metrics endpoint**: `http://localhost:8080/metrics`

Generate Swagger documentation:
```bash
make swagger
```

### Database Setup

The application uses MySQL with GORM for persistence. Configure using environment variables:

```bash
# Database configuration
export DATABASE_URL="root:password@tcp(localhost:3306)/protobuf_go?charset=utf8mb4&parseTime=True&loc=Local"
export DATABASE_MAX_IDLE=10
export DATABASE_MAX_OPEN=100
export DATABASE_MAX_LIFE=3600
```

**Docker MySQL Setup:**
```bash
docker run --name mysql-protobuf \
  -e MYSQL_ROOT_PASSWORD=password \
  -e MYSQL_DATABASE=protobuf_go \
  -p 3306:3306 \
  -d mysql:8.0
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
make swagger # Generate Swagger/OpenAPI documentation
make build   # Build static binary for production
make clean   # Remove generated files
make run     # Run development server
```

### Manual Proto Generation

If you prefer to generate manually:

```bash
make proto
```

## Validation & Database Models

The server uses **protoc-gen-validate** for automatic validation and **protoc-gen-gorm** for automatic GORM model generation from proto annotations. Both validation rules and database schema are defined directly in proto files.

### How It Works

1. **Define validation rules and GORM annotations in proto files**
2. **protoc-gen-validate generates validation code** automatically 
3. **protoc-gen-gorm generates GORM models** automatically
4. **Call `req.Validate()`** in service methods to validate requests
5. **Use generated ORM models** with `ToORM()` and `ToPB()` methods
6. **No manual validation or model code needed** - everything is generated from proto

### Validation Rules

**UserDTO Validation:**
```protobuf
message UserDTO {
  string name = 1 [(validate.rules).string = {min_len: 2, max_len: 100}];
  string email = 2 [(validate.rules).string = {pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$", max_len: 255}];
}
```

**CreateUserRequest Validation:**
```protobuf
message CreateUserRequest {
  UserDTO user = 1 [(validate.rules).message = {required: true}];
}
```

**GetUserRequest Validation:**
```protobuf
message GetUserRequest {
  int64 id = 1 [(validate.rules).int64 = {gt: 0}];
}
```

### GORM Model Annotations

**UserEntity with GORM Tags:**
```protobuf
message UserEntity {
  option (gorm.opts) = {
    ormable: true,
    table: "users"
  };
  
  uint32 id = 1 [(gorm.field).tag = {primary_key: true, auto_increment: true}];
  string name = 2 [(gorm.field).tag = {not_null: true, size: 100}];
  string email = 3 [(gorm.field).tag = {unique_index: "email_idx", size: 255}];
  string created_at = 4 [(gorm.field).tag = {not_null: true}];
  string updated_at = 5 [(gorm.field).tag = {not_null: true}];
}
```

### Generated Code Usage

**Validation in Service Methods:**
```go
func (s *UserServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
    // Validation is automatically generated from proto annotations
    if err := req.Validate(); err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
    }
    // ... rest of service logic
}
```

**GORM Model Usage:**
```go
// Create user entity from DTO
userEntity := &userpb.UserEntity{
    Name:  req.User.Name,
    Email: req.User.Email,
}

// Convert to ORM model for database operations
userORM, err := userEntity.ToORM(ctx)
if err != nil {
    return nil, status.Errorf(codes.Internal, "failed to process user data")
}

// Save to database
if err := s.userRepo.Create(ctx, &userORM); err != nil {
    return nil, status.Errorf(codes.Internal, "failed to create user")
}

// Convert back to protobuf for response
createdUser, err := userORM.ToPB(ctx)
if err != nil {
    return nil, status.Errorf(codes.Internal, "failed to process user data")
}
```

### Validation Examples

**Valid CreateUser Request:**
```json
{
  "user": {
    "name": "John Doe",
    "email": "john.doe@example.com"
  }
}
```

**Invalid Requests:**
```json
// Name too short (min_len: 2)
{"user": {"name": "J", "email": "john@example.com"}}
// Invalid email format (pattern validation)
{"user": {"name": "John Doe", "email": "invalid-email"}}
// Invalid user ID (gt: 0)
{"id": 0}
```

**Validation Error Response:**
```json
{
  "code": 3,
  "message": "validation failed: invalid CreateUserRequest.User: embedded message failed validation | caused by: invalid UserDTO.Name: value length must be at least 2 characters"
}
```

### Benefits

- **Single Source of Truth**: Validation rules and database schema defined in proto files
- **Automatic Code Generation**: No manual validation or model code needed
- **Type Safety**: Generated validation and models match proto definitions exactly
- **Consistent**: Same validation and data models across all languages
- **Maintainable**: Update proto annotations to change validation rules and database schema
- **Clean Architecture**: Repository pattern with generated GORM models
- **Seamless Conversion**: `ToORM()` and `ToPB()` methods for easy conversion

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
