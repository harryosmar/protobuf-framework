package usecase

import "errors"

// Usecase layer errors
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserEmailExists    = errors.New("user with email already exists")
	ErrInvalidUserData    = errors.New("invalid user data")
	ErrUserCreationFailed = errors.New("user creation failed")
	ErrUserUpdateFailed   = errors.New("user update failed")
	ErrUserDeletionFailed = errors.New("user deletion failed")
)
