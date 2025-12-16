package usecase

import (
	"context"
	"errors"
	"fmt"

	error2 "github.com/harryosmar/protobuf-go/error"
	userpb "github.com/harryosmar/protobuf-go/gen/user"
	"github.com/harryosmar/protobuf-go/repository"
)

// UserUsecase defines the interface for user business logic
type UserUsecase interface {
	CreateUser(ctx context.Context, userDTO *userpb.UserDTO) (*userpb.UserEntity, error)
	GetUserByID(ctx context.Context, id int64) (*userpb.UserEntity, error)
	GetUserByEmail(ctx context.Context, email string) (*userpb.UserEntity, error)
	UpdateUser(ctx context.Context, user *userpb.UserEntity) error
	DeleteUser(ctx context.Context, id int64) error
}

// userUsecase implements UserUsecase interface
type userUsecase struct {
	userRepo repository.UserRepository
}

// NewUserUsecase creates a new user usecase instance
func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

// CreateUser handles the business logic for creating a user
func (u *userUsecase) CreateUser(ctx context.Context, userDTO *userpb.UserDTO) (*userpb.UserEntity, error) {
	// Create user entity from DTO
	userEntity := &userpb.UserEntity{
		Name:      userDTO.Name,
		Email:     userDTO.Email,
		CreatedAt: "", // Will be set by database
		UpdatedAt: "", // Will be set by database
	}

	// Convert to ORM model for database operations
	userORM, err := userEntity.ToORM(ctx)
	if err != nil {
		return nil, error2.NewAppError(error2.ErrInvalidUserData, "failed to convert user entity to ORM", err)
	}

	// Save to database using repository
	if err := u.userRepo.Create(ctx, &userORM); err != nil {
		// Handle repository-specific errors and convert to usecase errors
		if errors.Is(err, repository.ErrUserEmailExists) {
			return nil, error2.NewAppError(error2.ErrUserEmailExists, fmt.Sprintf("user with email %s already exists", userDTO.Email), err)
		}
		return nil, error2.NewAppError(error2.ErrUserCreationFailed, "failed to create user in database", err)
	}

	// Convert back to protobuf entity for response
	createdUser, err := userORM.ToPB(ctx)
	if err != nil {
		return nil, error2.NewAppError(error2.ErrInvalidUserData, "failed to convert ORM to protobuf entity", err)
	}

	return &createdUser, nil
}

// GetUserByID handles the business logic for retrieving a user by ID
func (u *userUsecase) GetUserByID(ctx context.Context, id int64) (*userpb.UserEntity, error) {
	// Query database for user using repository
	userORM, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, error2.NewAppError(error2.ErrInternalServer, "failed to query user from database", err)
	}
	if userORM == nil {
		return nil, error2.NewAppError(error2.ErrUserNotFound, fmt.Sprintf("user with ID %d not found", id), nil)
	}

	// Convert ORM to protobuf entity
	user, err := userORM.ToPB(ctx)
	if err != nil {
		return nil, error2.NewAppError(error2.ErrInvalidUserData, "failed to convert ORM to protobuf entity", err)
	}

	return &user, nil
}

// GetUserByEmail handles the business logic for retrieving a user by email
func (u *userUsecase) GetUserByEmail(ctx context.Context, email string) (*userpb.UserEntity, error) {
	// Query database for user using repository
	userORM, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, error2.NewAppError(error2.ErrInternalServer, "failed to query user from database", err)
	}
	if userORM == nil {
		return nil, error2.NewAppError(error2.ErrUserNotFound, fmt.Sprintf("user with email %s not found", email), nil)
	}

	// Convert ORM to protobuf entity
	user, err := userORM.ToPB(ctx)
	if err != nil {
		return nil, error2.NewAppError(error2.ErrInvalidUserData, "failed to convert ORM to protobuf entity", err)
	}

	return &user, nil
}

// UpdateUser handles the business logic for updating a user
func (u *userUsecase) UpdateUser(ctx context.Context, user *userpb.UserEntity) error {
	// Convert to ORM model for database operations
	userORM, err := user.ToORM(ctx)
	if err != nil {
		return error2.NewAppError(error2.ErrInvalidUserData, "failed to convert user entity to ORM", err)
	}

	// Update in database using repository
	if err := u.userRepo.Update(ctx, &userORM); err != nil {
		return error2.NewAppError(error2.ErrUserUpdateFailed, "failed to update user in database", err)
	}
	return nil
}

// DeleteUser handles the business logic for deleting a user
func (u *userUsecase) DeleteUser(ctx context.Context, id int64) error {
	// Delete from database using repository
	if err := u.userRepo.Delete(ctx, id); err != nil {
		return error2.NewAppError(error2.ErrUserDeletionFailed, fmt.Sprintf("failed to delete user with ID %d", id), err)
	}
	return nil
}
