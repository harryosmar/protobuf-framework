package server

import (
	"context"

	appError "github.com/harryosmar/protobuf-go/error"
	userpb "github.com/harryosmar/protobuf-go/gen/user"
	"github.com/harryosmar/protobuf-go/logger"
	"github.com/harryosmar/protobuf-go/usecase"
	"go.uber.org/zap"
)

// UserServiceServer implements the UserService with usecase pattern
type UserServiceServer struct {
	userpb.UnimplementedUserServiceServer
	userUsecase usecase.UserServiceUsecase
}

// NewUserServiceServer creates a new UserServiceServer instance
func NewUserServiceServer(userUsecase usecase.UserServiceUsecase) *UserServiceServer {
	return &UserServiceServer{
		userUsecase: userUsecase,
	}
}

// CreateUser implements the CreateUser RPC method
func (s *UserServiceServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequestDTO) (*userpb.CreateUserResponseDTO, error) {
	var (
		log = logger.FromContext(ctx)
		err error
	)
	defer func() {
		if err != nil {
			log.Error("UserServiceServer.CreateUser err", zap.Error(err))
		}
	}()
	log.Info("UserService.CreateUser called", zap.String("name", req.User.Name), zap.String("email", req.User.Email))

	if err = req.Validate(); err != nil {
		return nil, appError.ErrInvalidArgument.WithMessage("validation failed: %v", err)
	}

	return s.userUsecase.CreateUser(ctx, req)
}

// GetUser implements the GetUser RPC method
func (s *UserServiceServer) GetUser(ctx context.Context, req *userpb.GetUserRequestDTO) (*userpb.GetUserResponse, error) {
	var (
		log = logger.FromContext(ctx)
		err error
	)
	defer func() {
		if err != nil {
			log.Error("UserServiceServer.GetUser err", zap.Error(err))
		}
	}()

	if err = req.Validate(); err != nil {
		return nil, appError.ErrInvalidArgument.WithMessage("validation failed: %v", err)
	}

	return s.userUsecase.GetUser(ctx, req)
}

// DeleteUser implements the DeleteUser RPC method
func (s *UserServiceServer) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequestDTO) (*userpb.DeleteUserResponseDTO, error) {
	var (
		log = logger.FromContext(ctx)
		err error
	)
	defer func() {
		if err != nil {
			log.Error("UserServiceServer.DeleteUser err", zap.Error(err))
		}
	}()

	if err = req.Validate(); err != nil {
		return nil, appError.ErrInvalidArgument.WithMessage("validation failed: %v", err)
	}

	return s.userUsecase.DeleteUser(ctx, req)
}

// UpdateUser implements the UpdateUser RPC method
func (s *UserServiceServer) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequestDTO) (*userpb.UpdateUserResponseDTO, error) {
	var (
		log = logger.FromContext(ctx)
		err error
	)
	defer func() {
		if err != nil {
			log.Error("UserServiceServer.UpdateUser err", zap.Error(err))
		}
	}()

	if err := req.Validate(); err != nil {
		return nil, appError.ErrInvalidArgument.WithMessage("validation failed: %v", err)
	}

	return s.userUsecase.UpdateUser(ctx, req)
}

// ListUsers implements the ListUsers RPC method
func (s *UserServiceServer) ListUsers(ctx context.Context, req *userpb.ListUsersRequestDTO) (*userpb.ListUsersResponseDTO, error) {
	var (
		log = logger.FromContext(ctx)
		err error
	)
	defer func() {
		if err != nil {
			log.Error("UserServiceServer.ListUser err", zap.Error(err))
		}
	}()

	if err := req.Validate(); err != nil {
		return nil, appError.ErrInvalidArgument.WithMessage("validation failed: %v", err)
	}

	return s.userUsecase.ListUsers(ctx, req)
}
