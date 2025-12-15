package service

import (
	"context"
	"errors"

	userpb "github.com/harryosmar/protobuf-go/gen/user"
	"github.com/harryosmar/protobuf-go/logger"
	"github.com/harryosmar/protobuf-go/repository"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserServer implements the UserService with repository pattern
type UserServer struct {
	userpb.UnimplementedUserServiceServer
	userRepo repository.UserRepository
}

// NewUserServer creates a new UserServer instance with dependency injection
func NewUserServer(userRepo repository.UserRepository) *UserServer {
	return &UserServer{
		userRepo: userRepo,
	}
}

// CreateUser implements the CreateUser RPC method
func (s *UserServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	// Get logger with request ID from context
	log := logger.FromContext(ctx)
	log.Info("UserService.CreateUser called", zap.String("name", req.User.Name), zap.String("email", req.User.Email))

	// Validation will be handled by protoc-gen-validate generated code
	// Proto validation rules:
	// - user: [(validate.rules).message = {required: true}]
	// - name: [(validate.rules).string = {min_len: 2, max_len: 100}]
	// - email: [(validate.rules).string = {pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$", max_len: 255}]
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	// Create user entity from DTO using generated GORM model
	userEntity := &userpb.UserEntity{
		Name:      req.User.Name,
		Email:     req.User.Email,
		CreatedAt: "", // Will be set by database
		UpdatedAt: "", // Will be set by database
	}

	// Convert to ORM model for database operations
	userORM, err := userEntity.ToORM(ctx)
	if err != nil {
		log.Error("Failed to convert user entity to ORM", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to process user data")
	}

	// Save to database using repository
	if err := s.userRepo.Create(ctx, &userORM); err != nil {
		log.Error("Failed to create user", zap.String("email", userORM.Email), zap.Error(err))

		// Handle repository-specific errors
		if errors.Is(err, repository.ErrUserEmailExists) {
			return nil, status.Errorf(codes.AlreadyExists, "user with email %s already exists", userORM.Email)
		}
		return nil, status.Errorf(codes.Internal, "failed to create user")
	}

	// Convert back to protobuf entity for response
	createdUser, err := userORM.ToPB(ctx)
	if err != nil {
		log.Error("Failed to convert ORM to protobuf", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to process user data")
	}

	log.Info("UserService.CreateUser created user", zap.Uint32("user_id", userORM.Id))
	return &userpb.CreateUserResponse{
		User: &createdUser,
	}, nil
}

// GetUser implements the GetUser RPC method
func (s *UserServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	// Get logger with request ID from context
	log := logger.FromContext(ctx)
	log.Info("UserService.GetUser called", zap.Int64("user_id", req.Id))

	// Validation will be handled by protoc-gen-validate generated code
	// Proto validation rule: [(validate.rules).int64 = {gt: 0}]
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	// Query database for user using repository
	userORM, err := s.userRepo.GetByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			log.Warn("UserService.GetUser user not found", zap.Int64("user_id", req.Id))
			return nil, status.Errorf(codes.NotFound, "user with ID %d not found", req.Id)
		}
		log.Error("Failed to query user", zap.Int64("user_id", req.Id), zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to retrieve user")
	}

	// Convert ORM to protobuf entity
	user, err := userORM.ToPB(ctx)
	if err != nil {
		log.Error("Failed to convert ORM to protobuf", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to process user data")
	}

	log.Info("UserService.GetUser found user", zap.String("user_name", user.Name))
	return &userpb.GetUserResponse{
		User: &user,
	}, nil
}
