package models

import (
	"time"

	userpb "github.com/harryosmar/protobuf-go/gen/user"
)

// User represents the user model for GORM
type User struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName returns the table name for the User model
func (User) TableName() string {
	return "users"
}

// ToProto converts the GORM User model to protobuf UserEntity
func (u *User) ToProto() *userpb.UserEntity {
	return &userpb.UserEntity{
		Id:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}

// FromProtoDTO creates a User model from protobuf UserDTO
func FromProtoDTO(dto *userpb.UserDTO) *User {
	return &User{
		Name:  dto.Name,
		Email: dto.Email,
	}
}
