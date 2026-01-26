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
	var (
		err error
		log = logger.FromContext(ctx)
	)
	log.Info("UserService.CreateUser called", zap.String("name", req.User.Name), zap.String("email", req.User.Email))

	if err = req.Validate(); err != nil {
		return nil, appError.ErrInvalidArgument.WithMessage("validation failed: %v", err)
	}

	createdUser, err := s.userUsecase.CreateUser(ctx, req.User)
	if err != nil {
		return nil, err
	}

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
		return nil, appError.ErrInvalidArgument.WithMessage("validation failed: %v", err)
	}

	// Call usecase to handle business logic
	user, err := s.userUsecase.GetUserByID(ctx, req.Id)
	if err != nil {
		log.Error("Failed to get user", zap.Int64("user_id", req.Id), zap.Error(err))
		// Error conversion handled automatically by ErrorConversionInterceptor
		return nil, err
	}

	log.Info("UserService.GetUser found user", zap.String("user_name", user.Name))
	return &userpb.GetUserResponse{
		User: user,
	}, nil
}
