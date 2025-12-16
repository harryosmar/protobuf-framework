package repository

import (
	"context"
	userpb "github.com/harryosmar/protobuf-go/gen/user"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *userpb.UserEntityORM) error
	GetByID(ctx context.Context, id int64) (*userpb.UserEntityORM, error)
	GetByEmail(ctx context.Context, email string) (*userpb.UserEntityORM, error)
	Update(ctx context.Context, user *userpb.UserEntityORM) error
	Delete(ctx context.Context, id int64) error
}
