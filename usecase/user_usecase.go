package usecase

import (
	"context"

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
		return nil, err
	}

	// Save to database using repository
	if err := u.userRepo.Create(ctx, &userORM); err != nil {
		return nil, err
	}

	// Convert back to protobuf entity for response
	createdUser, err := userORM.ToPB(ctx)
	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}

// GetUserByID handles the business logic for retrieving a user by ID
func (u *userUsecase) GetUserByID(ctx context.Context, id int64) (*userpb.UserEntity, error) {
	// Query database for user using repository
	userORM, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if userORM == nil {
		return nil, error2.ErrUserNotFound.WithMessage("user with ID %d not found", id)
	}

	// Convert ORM to protobuf entity
	user, err := userORM.ToPB(ctx)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail handles the business logic for retrieving a user by email
func (u *userUsecase) GetUserByEmail(ctx context.Context, email string) (*userpb.UserEntity, error) {
	// Query database for user using repository
	userORM, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if userORM == nil {
		return nil, error2.ErrUserNotFound.WithMessage("user with email %s not found", email)
	}

	// Convert ORM to protobuf entity
	user, err := userORM.ToPB(ctx)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser handles the business logic for updating a user
func (u *userUsecase) UpdateUser(ctx context.Context, user *userpb.UserEntity) error {
	// Convert to ORM model for database operations
	userORM, err := user.ToORM(ctx)
	if err != nil {
		return err
	}

	// Update in database using repository
	return u.userRepo.Update(ctx, &userORM)
}

// DeleteUser handles the business logic for deleting a user
func (u *userUsecase) DeleteUser(ctx context.Context, id int64) error {
	// Delete from database using repository
	return u.userRepo.Delete(ctx, id)
}
