package error

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
	Code     string
	Status   int
	Message  string
	GrpcCode codes.Code
}

var (
	internalServerError = CodeErrEntity{Code: "ERR500P00", Status: http.StatusInternalServerError, GrpcCode: codes.Internal, Message: "internal server error"}
)

// Error codes following pattern: ERRXXXPYY
// XXX: HTTP status code (400, 404, 409, 500, etc.)
// P: Identifier for protobuf-go service
// YY: Incremental error number starting from 00
var (
	codeErrMap = map[CodeErr]CodeErrEntity{
		ErrInternalServer:     internalServerError,
		ErrValidationFailed:   {Code: "ERR400P01", Status: http.StatusBadRequest, GrpcCode: codes.InvalidArgument, Message: "validation failed"},
		ErrNotFound:           {Code: "ERR404P02", Status: http.StatusNotFound, GrpcCode: codes.NotFound, Message: "not found"},
		ErrUnauthorized:       {Code: "ERR401P03", Status: http.StatusUnauthorized, GrpcCode: codes.Unauthenticated, Message: "unauthorized"},
		ErrForbidden:          {Code: "ERR403P04", Status: http.StatusForbidden, GrpcCode: codes.PermissionDenied, Message: "forbidden"},
		ErrTooManyRequest:     {Code: "ERR429P05", Status: http.StatusTooManyRequests, GrpcCode: codes.ResourceExhausted, Message: "too many requests"},
		ErrUserNotFound:       {Code: "ERR404P06", Status: http.StatusNotFound, GrpcCode: codes.NotFound, Message: "user not found"},
		ErrUserEmailExists:    {Code: "ERR409P07", Status: http.StatusConflict, GrpcCode: codes.AlreadyExists, Message: "user with email already exists"},
		ErrInvalidUserData:    {Code: "ERR400P07", Status: http.StatusBadRequest, GrpcCode: codes.InvalidArgument, Message: "invalid user data"},
		ErrUserCreationFailed: {Code: "ERR500P09", Status: http.StatusInternalServerError, GrpcCode: codes.Internal, Message: "user creation failed"},
		ErrUserUpdateFailed:   {Code: "ERR500P10", Status: http.StatusInternalServerError, GrpcCode: codes.Internal, Message: "user update failed"},
		ErrUserDeletionFailed: {Code: "ERR500P11", Status: http.StatusInternalServerError, GrpcCode: codes.Internal, Message: "user deletion failed"},
	}
)

// Error code constants
const (
	ErrInternalServer CodeErr = iota + 100
	ErrValidationFailed
	ErrNotFound
	ErrUnauthorized
	ErrForbidden
	ErrTooManyRequest
	ErrUserNotFound
	ErrUserEmailExists
	ErrInvalidUserData
	ErrUserCreationFailed
	ErrUserUpdateFailed
	ErrUserDeletionFailed
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

// GetCodeErrEntity returns the error code entity
func (e *AppError) GetCodeErrEntity() CodeErrEntity {
	if entity, exists := codeErrMap[e.Code]; exists {
		return entity
	}
	return internalServerError
}

// GetCode returns the error code
func (e *AppError) GetCode() string {
	return e.GetCodeErrEntity().Code
}

// GetStatus returns the HTTP status code
func (e *AppError) GetStatus() int {
	return e.GetCodeErrEntity().Status
}

// GetMessage returns the error message
func (e *AppError) GetMessage() string {
	return e.GetCodeErrEntity().Message
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
	return status.Error(e.GetCodeErrEntity().GrpcCode, e.GetMessage())
}

// IsErrorCode checks if an error is of a specific code
func IsErrorCode(err error, code CodeErr) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == code
	}
	return false
}
