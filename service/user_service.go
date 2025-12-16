package service

import (
	"context"
	"errors"

	userpb "github.com/harryosmar/protobuf-go/gen/user"
	"github.com/harryosmar/protobuf-go/logger"
	"github.com/harryosmar/protobuf-go/usecase"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserServiceServer implements the UserService with usecase pattern
type UserServiceServer struct {
	userpb.UnimplementedUserServiceServer
	userUsecase usecase.UserUsecase
}

// NewUserServiceServer creates a new UserServiceServer instance
func NewUserServiceServer(userUsecase usecase.UserUsecase) *UserServiceServer {
	return &UserServiceServer{
		userUsecase: userUsecase,
	}
}

// CreateUser implements the CreateUser RPC method
func (s *UserServiceServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	// Get logger with request ID from context
	log := logger.FromContext(ctx)
	log.Info("UserService.CreateUser called", zap.String("name", req.User.Name), zap.String("email", req.User.Email))

	// Validation will be handled by protoc-gen-validate generated code
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	// Call usecase to handle business logic
	createdUser, err := s.userUsecase.CreateUser(ctx, req.User)
	if err != nil {
		log.Error("Failed to create user", zap.String("email", req.User.Email), zap.Error(err))

		// Handle usecase-specific errors and map to gRPC status codes
		if errors.Is(err, usecase.ErrUserEmailExists) {
			return nil, status.Errorf(codes.AlreadyExists, "user with email %s already exists", req.User.Email)
		}
		if errors.Is(err, usecase.ErrInvalidUserData) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid user data")
		}
		return nil, status.Errorf(codes.Internal, "failed to create user")
	}

	log.Info("UserService.CreateUser created user", zap.String("user_name", createdUser.Name))
	return &userpb.CreateUserResponse{
		User: createdUser,
	}, nil
}

// GetUser implements the GetUser RPC method
func (s *UserServiceServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	// Get logger with request ID from context
	log := logger.FromContext(ctx)
	log.Info("UserService.GetUser called", zap.Int64("user_id", req.Id))

	// Validation will be handled by protoc-gen-validate generated code
	// Proto validation rule: [(validate.rules).int64 = {gt: 0}]
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	// Call usecase to handle business logic
	user, err := s.userUsecase.GetUserByID(ctx, req.Id)
	if err != nil {
		log.Error("Failed to get user", zap.Int64("user_id", req.Id), zap.Error(err))

		// Handle usecase-specific errors and map to gRPC status codes
		if errors.Is(err, usecase.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user with ID %d not found", req.Id)
		}
		if errors.Is(err, usecase.ErrInvalidUserData) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid user data")
		}
		return nil, status.Errorf(codes.Internal, "failed to retrieve user")
	}

	log.Info("UserService.GetUser found user", zap.String("user_name", user.Name))
	return &userpb.GetUserResponse{
		User: user,
	}, nil
}
