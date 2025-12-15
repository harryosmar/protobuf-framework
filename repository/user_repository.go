package repository

import (
	"context"
	"strings"

	userpb "github.com/harryosmar/protobuf-go/gen/user"
	"gorm.io/gorm"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *userpb.UserEntityORM) error
	GetByID(ctx context.Context, id int64) (*userpb.UserEntityORM, error)
	GetByEmail(ctx context.Context, email string) (*userpb.UserEntityORM, error)
	Update(ctx context.Context, user *userpb.UserEntityORM) error
	Delete(ctx context.Context, id int64) error
}

// userRepository implements UserRepository interface
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create creates a new user in the database
func (r *userRepository) Create(ctx context.Context, user *userpb.UserEntityORM) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		// Check for duplicate email error
		if strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "UNIQUE constraint") {
			return ErrUserEmailExists
		}
		return err
	}
	return nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id int64) (*userpb.UserEntityORM, error) {
	var user userpb.UserEntityORM
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*userpb.UserEntityORM, error) {
	var user userpb.UserEntityORM
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// Update updates an existing user
func (r *userRepository) Update(ctx context.Context, user *userpb.UserEntityORM) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete deletes a user by ID
func (r *userRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&userpb.UserEntityORM{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}
