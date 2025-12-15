package service

import (
	"context"
	"errors"
	"strings"

	userpb "github.com/harryosmar/protobuf-go/gen/user"
	"github.com/harryosmar/protobuf-go/logger"
	"github.com/harryosmar/protobuf-go/models"
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

	// Input validation
	if err := s.validateCreateUserRequest(req); err != nil {
		return nil, err
	}

	// Create user model from DTO
	user := models.FromProtoDTO(req.User)

	// Save to database using repository
	if err := s.userRepo.Create(ctx, user); err != nil {
		log.Error("Failed to create user", zap.String("email", user.Email), zap.Error(err))

		// Handle repository-specific errors
		if errors.Is(err, repository.ErrUserEmailExists) {
			return nil, status.Errorf(codes.AlreadyExists, "user with email %s already exists", user.Email)
		}
		return nil, status.Errorf(codes.Internal, "failed to create user")
	}

	log.Info("UserService.CreateUser created user", zap.Int64("user_id", user.ID))
	return &userpb.CreateUserResponse{
		User: user.ToProto(),
	}, nil
}

// validateCreateUserRequest validates the create user request
func (s *UserServer) validateCreateUserRequest(req *userpb.CreateUserRequest) error {
	if req.User == nil {
		return status.Errorf(codes.InvalidArgument, "user data is required")
	}
	if strings.TrimSpace(req.User.Name) == "" {
		return status.Errorf(codes.InvalidArgument, "user name is required")
	}
	if strings.TrimSpace(req.User.Email) == "" {
		return status.Errorf(codes.InvalidArgument, "user email is required")
	}
	return nil
}

// GetUser implements the GetUser RPC method
func (s *UserServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	// Get logger with request ID from context
	log := logger.FromContext(ctx)
	log.Info("UserService.GetUser called", zap.Int64("user_id", req.Id))

	// Input validation
	if req.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "user ID must be positive")
	}

	// Query database for user using repository
	user, err := s.userRepo.GetByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			log.Warn("UserService.GetUser user not found", zap.Int64("user_id", req.Id))
			return nil, status.Errorf(codes.NotFound, "user with ID %d not found", req.Id)
		}
		log.Error("Failed to query user", zap.Int64("user_id", req.Id), zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to retrieve user")
	}

	log.Info("UserService.GetUser found user", zap.String("user_name", user.Name))
	return &userpb.GetUserResponse{
		User: user.ToProto(),
	}, nil
}
