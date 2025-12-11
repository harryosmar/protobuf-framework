package service

import (
	"context"
	"fmt"
	"time"

	userpb "github.com/harryosmar/protobuf-go/gen/user"
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

	return &userpb.CreateUserResponse{
		User: user,
	}, nil
}

// GetUser implements the GetUser RPC method
func (s *UserServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	user, exists := s.users[req.Id]
	if !exists {
		return nil, fmt.Errorf("user with ID %d not found", req.Id)
	}

	return &userpb.GetUserResponse{
		User: user,
	}, nil
}
