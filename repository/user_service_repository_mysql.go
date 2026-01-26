package repository

import (
	"context"
	"errors"

	appErrors "github.com/harryosmar/protobuf-go/error"

	"github.com/go-sql-driver/mysql"
	userpb "github.com/harryosmar/protobuf-go/gen/user"
	"gorm.io/gorm"
)

// userServiceRepositoryMySQL implements UserServiceRepository interface
type userServiceRepositoryMySQL struct {
	db *gorm.DB
}

// NewUserServiceRepositoryMySQL creates a new user repository instance
func NewUserServiceRepositoryMySQL(db *gorm.DB) UserServiceRepository {
	return &userServiceRepositoryMySQL{
		db: db,
	}
}

// Create creates a new user in the database
func (r *userServiceRepositoryMySQL) Create(ctx context.Context, user *userpb.UserEntityORM) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		// Check for MySQL duplicate entry error (Error 1062)
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return appErrors.ErrUserEmailExists
		}
		return err
	}
	return nil
}

// GetByID retrieves a user by ID
func (r *userServiceRepositoryMySQL) GetByID(ctx context.Context, id int64) (*userpb.UserEntityORM, error) {
	var user userpb.UserEntityORM
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found is not an error at repository level
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *userServiceRepositoryMySQL) GetByEmail(ctx context.Context, email string) (*userpb.UserEntityORM, error) {
	var user userpb.UserEntityORM
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found is not an error at repository level
		}
		return nil, err
	}
	return &user, nil
}

// Update updates an existing user
func (r *userServiceRepositoryMySQL) Update(ctx context.Context, user *userpb.UserEntityORM) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete deletes a user by ID
func (r *userServiceRepositoryMySQL) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&userpb.UserEntityORM{}, id)
	if result.Error != nil {
		return result.Error
	}
	// Return success even if no rows affected - idempotent delete
	return nil
}
