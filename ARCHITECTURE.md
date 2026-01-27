# Protobuf-Go Architecture

This document describes the standardized architecture of the protobuf-go project, focusing on the server, usecase, and repository layers.

## Project Structure

```
protobuf-go/
├── proto/              # Protocol buffer definitions
│   ├── hello.proto
│   └── user.proto
├── gen/                # Generated code from proto files
│   ├── hello/
│   └── user/
├── server/             # gRPC server implementations
│   ├── hello_service_server.go
│   ├── user_service_server.go
│   └── user_service_server_additional.go
├── usecase/            # Business logic layer
│   ├── hello_service_usecase.go
│   ├── user_service_usecase.go
│   └── user_service_handler_methods.go
├── repository/         # Data access layer
│   ├── user_service_repository.go
│   └── user_service_repository_mysql.go
├── error/              # Error handling
│   └── codes.go
├── middleware/         # gRPC middleware
│   └── error_conversion.go
├── config/             # Application configuration
│   └── config.go
├── database/           # Database configuration
│   └── database.go
└── main.go             # Application entry point
```

## Architecture Overview

The project follows a clean architecture pattern with clear separation of concerns:

1. **Server Layer**: Handles gRPC requests, validation, and delegates to the usecase layer
2. **Usecase Layer**: Contains business logic and orchestrates repository calls
3. **Repository Layer**: Manages data access and database operations

## Interface Standardization

### Server Layer

The server layer implements the gRPC service interfaces generated from proto files:

```go
// UserServiceServer implements the UserService with usecase pattern
type UserServiceServer struct {
    userpb.UnimplementedUserServiceServer
    userUsecase usecase.UserServiceUsecase
}
```

Server methods follow a consistent pattern:
1. Extract logger from context
2. Log the request
3. Validate the request using generated validation code
4. Call the appropriate usecase method
5. Return the response or error

### Usecase Layer

The usecase layer defines interfaces for business operations:

```go
// UserServiceUsecase defines the interface for user business logic
type UserServiceUsecase interface {
    // Core business logic methods
    CreateUserEntity(ctx context.Context, userDTO *userpb.UserDTO) (*userpb.UserEntity, error)
    GetUserByID(ctx context.Context, id int64) (*userpb.UserEntity, error)
    GetUserByEmail(ctx context.Context, email string) (*userpb.UserEntity, error)
    UpdateUserEntity(ctx context.Context, user *userpb.UserEntity) error
    DeleteUserByID(ctx context.Context, id int64) error
    
    // Handler interface methods
    CreateUser(ctx context.Context, req *userpb.CreateUserRequestDTO) (*userpb.CreateUserResponseDTO, error)
    GetUser(ctx context.Context, req *userpb.GetUserRequestDTO) (*userpb.GetUserResponse, error)
    UpdateUser(ctx context.Context, req *userpb.UpdateUserRequestDTO) (*userpb.UpdateUserResponseDTO, error)
    DeleteUser(ctx context.Context, req *userpb.DeleteUserRequestDTO) (*userpb.DeleteUserResponseDTO, error)
    ListUsers(ctx context.Context, req *userpb.ListUsersRequestDTO) (*userpb.ListUsersResponseDTO, error)
}
```

The usecase implementation is divided into two parts:
1. Core business logic methods that handle domain operations
2. Handler interface methods that adapt between gRPC DTOs and domain entities

### Repository Layer

The repository layer defines interfaces for data operations:

```go
// UserServiceRepository defines the interface for user data operations
type UserServiceRepository interface {
    Create(ctx context.Context, user *userpb.UserEntityORM) error
    GetByID(ctx context.Context, id int64) (*userpb.UserEntityORM, error)
    GetByEmail(ctx context.Context, email string) (*userpb.UserEntityORM, error)
    Update(ctx context.Context, user *userpb.UserEntityORM) error
    Delete(ctx context.Context, id int64) error
}
```

The repository implementation handles database operations using GORM.

## Error Handling

The project uses a standardized error handling approach:

1. **Domain Errors**: Defined in the `error` package with specific error codes
2. **Error Conversion**: Middleware automatically converts domain errors to gRPC status codes
3. **Consistent Error Format**: All errors follow a consistent format with code, message, and HTTP status

## Dependency Injection

The project uses constructor-based dependency injection:

```go
// Initialize repositories
userRepo := repository.NewUserServiceRepositoryMySQL(db)

// Initialize usecases
userUsecase := usecase.NewUserServiceUsecase(userRepo)

// Initialize servers
userServer := server.NewUserServiceServer(userUsecase)
```

This approach ensures that dependencies are explicitly declared and makes testing easier.

## Code Generation

The project uses several code generators:

1. **protoc-gen-go**: Generates Go code from proto files
2. **protoc-gen-go-grpc**: Generates gRPC service code
3. **protoc-gen-grpc-gateway**: Generates HTTP gateway code
4. **protoc-gen-validate**: Generates validation code
5. **protoc-gen-gorm**: Generates GORM models

## Best Practices

1. **Interface Segregation**: Separate interfaces for different layers
2. **Single Responsibility**: Each component has a clear responsibility
3. **Dependency Inversion**: High-level modules depend on abstractions
4. **Error Handling**: Consistent error handling across all layers
5. **Logging**: Structured logging with request IDs for traceability
6. **Validation**: Automatic validation from proto annotations
7. **Database Access**: Repository pattern for data access abstraction
