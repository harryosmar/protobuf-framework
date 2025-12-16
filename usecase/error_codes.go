package usecase

import (
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CodeErr represents an error code type
type CodeErr int

// CodeErrEntity represents the structure of an error with code, status, and message
type CodeErrEntity struct {
	Code    string
	Status  int
	Message string
}

// Error codes following pattern: ERRXXXPYY
// XXX: HTTP status code (400, 404, 409, 500, etc.)
// P: Identifier for protobuf-go service
// YY: Incremental error number starting from 00
var (
	codeErrMap = map[CodeErr]CodeErrEntity{
		ErrUserNotFound:       {Code: "ERR404P00", Status: http.StatusNotFound, Message: "user not found"},
		ErrUserEmailExists:    {Code: "ERR409P01", Status: http.StatusConflict, Message: "user with email already exists"},
		ErrInvalidUserData:    {Code: "ERR400P02", Status: http.StatusBadRequest, Message: "invalid user data"},
		ErrUserCreationFailed: {Code: "ERR500P03", Status: http.StatusInternalServerError, Message: "user creation failed"},
		ErrUserUpdateFailed:   {Code: "ERR500P04", Status: http.StatusInternalServerError, Message: "user update failed"},
		ErrUserDeletionFailed: {Code: "ERR500P05", Status: http.StatusInternalServerError, Message: "user deletion failed"},
		ErrValidationFailed:   {Code: "ERR400P06", Status: http.StatusBadRequest, Message: "validation failed"},
		ErrInternalServer:     {Code: "ERR500P07", Status: http.StatusInternalServerError, Message: "internal server error"},
	}
)

// Error code constants
const (
	ErrUserNotFound CodeErr = iota + 100
	ErrUserEmailExists
	ErrInvalidUserData
	ErrUserCreationFailed
	ErrUserUpdateFailed
	ErrUserDeletionFailed
	ErrValidationFailed
	ErrInternalServer
)

// AppError represents a structured application error
type AppError struct {
	Code    CodeErr
	Message string
	Err     error
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// GetCode returns the error code
func (e *AppError) GetCode() string {
	if entity, exists := codeErrMap[e.Code]; exists {
		return entity.Code
	}
	return "ERR500P99" // Unknown error
}

// GetStatus returns the HTTP status code
func (e *AppError) GetStatus() int {
	if entity, exists := codeErrMap[e.Code]; exists {
		return entity.Status
	}
	return http.StatusInternalServerError
}

// GetMessage returns the error message
func (e *AppError) GetMessage() string {
	if e.Message != "" {
		return e.Message
	}
	if entity, exists := codeErrMap[e.Code]; exists {
		return entity.Message
	}
	return "unknown error"
}

// NewAppError creates a new application error
func NewAppError(code CodeErr, message string, err error) *AppError {
	if message == "" {
		if entity, exists := codeErrMap[code]; exists {
			message = entity.Message
		}
	}
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// ToGRPCStatus converts AppError to gRPC status
func (e *AppError) ToGRPCStatus() error {
	var code codes.Code
	switch e.Code {
	case ErrUserNotFound:
		code = codes.NotFound
	case ErrUserEmailExists:
		code = codes.AlreadyExists
	case ErrInvalidUserData, ErrValidationFailed:
		code = codes.InvalidArgument
	case ErrUserCreationFailed, ErrUserUpdateFailed, ErrUserDeletionFailed, ErrInternalServer:
		code = codes.Internal
	default:
		code = codes.Internal
	}
	return status.Error(code, e.GetMessage())
}

// IsErrorCode checks if an error is of a specific code
func IsErrorCode(err error, code CodeErr) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == code
	}
	return false
}
