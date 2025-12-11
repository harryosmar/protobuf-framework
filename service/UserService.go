package service

import (
	"context"
	"fmt"
	"time"

	userpb "github.com/harryosmar/protobuf-go/gen/user"
	"github.com/harryosmar/protobuf-go/logger"
	"go.uber.org/zap"
)

// UserServer implements the UserService
type UserServer struct {
	userpb.UnimplementedUserServiceServer
	users  map[int64]*userpb.UserEntity
	nextID int64
}

// NewUserServer creates a new UserServer instance
func NewUserServer() *UserServer {
	return &UserServer{
		users:  make(map[int64]*userpb.UserEntity),
		nextID: 1,
	}
}

// CreateUser implements the CreateUser RPC method
func (s *UserServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	// Get logger with request ID from context
	log := logger.FromContext(ctx)
	log.Info("UserService.CreateUser called", zap.String("name", req.User.Name), zap.String("email", req.User.Email))

	now := time.Now().Format(time.RFC3339)

	user := &userpb.UserEntity{
		Id:        s.nextID,
		Name:      req.User.Name,
		Email:     req.User.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.users[s.nextID] = user
	s.nextID++

	log.Info("UserService.CreateUser created user", zap.Int64("user_id", user.Id))
	return &userpb.CreateUserResponse{
		User: user,
	}, nil
}

// GetUser implements the GetUser RPC method
func (s *UserServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	// Get logger with request ID from context
	log := logger.FromContext(ctx)
	log.Info("UserService.GetUser called", zap.Int64("user_id", req.Id))

	user, exists := s.users[req.Id]
	if !exists {
		log.Warn("UserService.GetUser user not found", zap.Int64("user_id", req.Id))
		return nil, fmt.Errorf("user with ID %d not found", req.Id)
	}

	log.Info("UserService.GetUser found user", zap.String("user_name", user.Name))
	return &userpb.GetUserResponse{
		User: user,
	}, nil
}
