# Protobuf Go gRPC Server

A simple gRPC server with HTTP gateway demonstrating Protocol Buffers and gRPC in Go.

## Features

- **gRPC Service**: `HelloService` with `GetHello` RPC method
- **HTTP Gateway**: REST API endpoint using grpc-gateway
- **Protocol Buffers**: Message definitions with `HelloRequest` and `HelloResponse`

## Project Structure

```
.
├── proto/              # Protocol buffer definitions
│   └── hello.proto
├── gen/                # Generated code from proto files
│   ├── hello.pb.go
│   ├── hello_grpc.pb.go
│   └── hello.pb.gw.go
├── server/             # Server implementation
│   └── server.go
├── main.go             # Main application
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

### Start the server

```bash
go run main.go
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

Using curl:

```bash
curl http://localhost:8080/v1/hello/World
```

Expected response:

```json
{"message":"Hello, World!"}
```

## Regenerating Proto Files

If you modify `proto/hello.proto`, regenerate the Go code:

```bash
protoc -I./proto -I$HOME/.proto \
  --go_out=./gen --go_opt=paths=source_relative \
  --go-grpc_out=./gen --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=./gen --grpc-gateway_opt=paths=source_relative \
  proto/hello.proto
```

## API Documentation

### GetHello

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
