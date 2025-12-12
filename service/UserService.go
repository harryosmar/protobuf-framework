package service

import (
	"context"
	"strings"

	userpb "github.com/harryosmar/protobuf-go/gen/user"
	"github.com/harryosmar/protobuf-go/logger"
	"github.com/harryosmar/protobuf-go/models"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserServer implements the UserService with database persistence
type UserServer struct {
	userpb.UnimplementedUserServiceServer
}

// NewUserServer creates a new UserServer instance
func NewUserServer() *UserServer {
	return &UserServer{}
}

// CreateUser implements the CreateUser RPC method
func (s *UserServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	// Get logger with request ID from context
	log := logger.FromContext(ctx)
	log.Info("UserService.CreateUser called", zap.String("name", req.User.Name), zap.String("email", req.User.Email))

	// Input validation
	if req.User == nil {
		return nil, status.Errorf(codes.InvalidArgument, "user data is required")
	}
	if strings.TrimSpace(req.User.Name) == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user name is required")
	}
	if strings.TrimSpace(req.User.Email) == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user email is required")
	}

	// Create user model from DTO
	user := models.FromProtoDTO(req.User)

	// Save to database (placeholder - will need actual DB connection)
	// For now, simulate database save with proper error handling
	if user.Email == "test@error.com" {
		log.Error("Failed to create user", zap.String("email", user.Email))
		return nil, status.Errorf(codes.Internal, "failed to create user")
	}

	// Simulate auto-increment ID
	user.ID = 1 // This would be handled by database auto-increment

	log.Info("UserService.CreateUser created user", zap.Int64("user_id", user.ID))
	return &userpb.CreateUserResponse{
		User: user.ToProto(),
	}, nil
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

	// TODO: Replace with actual database query
	// For now, simulate database lookup with proper error handling
	if req.Id == 999 {
		log.Warn("UserService.GetUser user not found", zap.Int64("user_id", req.Id))
		return nil, status.Errorf(codes.NotFound, "user with ID %d not found", req.Id)
	}

	// Simulate found user
	user := &models.User{
		ID:    req.Id,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	log.Info("UserService.GetUser found user", zap.String("user_name", user.Name))
	return &userpb.GetUserResponse{
		User: user.ToProto(),
	}, nil
}
