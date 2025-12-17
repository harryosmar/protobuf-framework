package error

import (
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CodeErr represents an error code type that implements error interface
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
		// Standard gRPC status codes
		ErrInternalServer:     internalServerError,
		ErrCancelled:          {Code: "ERR499P01", Status: 499, GrpcCode: codes.Canceled, Message: "request cancelled"},
		ErrUnknown:            {Code: "ERR500P02", Status: http.StatusInternalServerError, GrpcCode: codes.Unknown, Message: "unknown error"},
		ErrInvalidArgument:    {Code: "ERR400P03", Status: http.StatusBadRequest, GrpcCode: codes.InvalidArgument, Message: "invalid argument"},
		ErrDeadlineExceeded:   {Code: "ERR504P04", Status: http.StatusGatewayTimeout, GrpcCode: codes.DeadlineExceeded, Message: "deadline exceeded"},
		ErrNotFound:           {Code: "ERR404P05", Status: http.StatusNotFound, GrpcCode: codes.NotFound, Message: "not found"},
		ErrAlreadyExists:      {Code: "ERR409P06", Status: http.StatusConflict, GrpcCode: codes.AlreadyExists, Message: "already exists"},
		ErrPermissionDenied:   {Code: "ERR403P07", Status: http.StatusForbidden, GrpcCode: codes.PermissionDenied, Message: "permission denied"},
		ErrResourceExhausted:  {Code: "ERR429P08", Status: http.StatusTooManyRequests, GrpcCode: codes.ResourceExhausted, Message: "resource exhausted"},
		ErrFailedPrecondition: {Code: "ERR400P09", Status: http.StatusBadRequest, GrpcCode: codes.FailedPrecondition, Message: "failed precondition"},
		ErrAborted:            {Code: "ERR409P10", Status: http.StatusConflict, GrpcCode: codes.Aborted, Message: "aborted"},
		ErrOutOfRange:         {Code: "ERR400P11", Status: http.StatusBadRequest, GrpcCode: codes.OutOfRange, Message: "out of range"},
		ErrUnimplemented:      {Code: "ERR501P12", Status: http.StatusNotImplemented, GrpcCode: codes.Unimplemented, Message: "unimplemented"},
		ErrUnavailable:        {Code: "ERR503P14", Status: http.StatusServiceUnavailable, GrpcCode: codes.Unavailable, Message: "service unavailable"},
		ErrDataLoss:           {Code: "ERR500P15", Status: http.StatusInternalServerError, GrpcCode: codes.DataLoss, Message: "data loss"},
		ErrUnauthenticated:    {Code: "ERR401P16", Status: http.StatusUnauthorized, GrpcCode: codes.Unauthenticated, Message: "unauthenticated"},

		// Application-specific errors
		ErrUserNotFound:       {Code: "ERR404P17", Status: http.StatusNotFound, GrpcCode: codes.NotFound, Message: "user not found"},
		ErrUserEmailExists:    {Code: "ERR409P18", Status: http.StatusConflict, GrpcCode: codes.AlreadyExists, Message: "user with email already exists"},
		ErrInvalidUserData:    {Code: "ERR400P19", Status: http.StatusBadRequest, GrpcCode: codes.InvalidArgument, Message: "invalid user data"},
		ErrUserCreationFailed: {Code: "ERR500P20", Status: http.StatusInternalServerError, GrpcCode: codes.Internal, Message: "user creation failed"},
		ErrUserUpdateFailed:   {Code: "ERR500P21", Status: http.StatusInternalServerError, GrpcCode: codes.Internal, Message: "user update failed"},
		ErrUserDeletionFailed: {Code: "ERR500P22", Status: http.StatusInternalServerError, GrpcCode: codes.Internal, Message: "user deletion failed"},
	}
)

// Error code constants
const (
	// Standard gRPC status codes (following gRPC specification)
	ErrOK CodeErr = iota + 100
	ErrCancelled
	ErrUnknown
	ErrInvalidArgument
	ErrDeadlineExceeded
	ErrNotFound
	ErrAlreadyExists
	ErrPermissionDenied
	ErrResourceExhausted
	ErrFailedPrecondition
	ErrAborted
	ErrOutOfRange
	ErrUnimplemented
	ErrInternalServer
	ErrUnavailable
	ErrDataLoss
	ErrUnauthenticated

	// Application-specific errors
	ErrUserNotFound
	ErrUserEmailExists
	ErrInvalidUserData
	ErrUserCreationFailed
	ErrUserUpdateFailed
	ErrUserDeletionFailed
)

// Error implements the error interface for CodeErr
func (c CodeErr) Error() string {
	if entity, exists := codeErrMap[c]; exists {
		return entity.Message
	}
	return internalServerError.Message
}

// GetCodeErrEntity returns the error code entity
func (c CodeErr) GetCodeErrEntity() CodeErrEntity {
	if entity, exists := codeErrMap[c]; exists {
		return entity
	}
	return internalServerError
}

// GetCode returns the error code
func (c CodeErr) GetCode() string {
	return c.GetCodeErrEntity().Code
}

// GetStatus returns the HTTP status code
func (c CodeErr) GetStatus() int {
	return c.GetCodeErrEntity().Status
}

// GetMessage returns the error message
func (c CodeErr) GetMessage() string {
	return c.GetCodeErrEntity().Message
}

// ToGRPCStatus converts CodeErr to gRPC status
func (c CodeErr) ToGRPCStatus() error {
	return status.Error(c.GetCodeErrEntity().GrpcCode, c.GetMessage())
}

// WithMessage returns a formatted error with additional context while preserving CodeErr type
func (c CodeErr) WithMessage(format string, args ...interface{}) *CodeErrWithContext {
	baseMessage := c.GetMessage()
	if format != "" {
		customMessage := fmt.Sprintf(format, args...)
		return &CodeErrWithContext{
			CodeErr: c,
			message: fmt.Sprintf("%s: %s", baseMessage, customMessage),
		}
	}
	return &CodeErrWithContext{
		CodeErr: c,
		message: baseMessage,
	}
}

// CodeErrWithContext wraps CodeErr with additional context while preserving gRPC compatibility
type CodeErrWithContext struct {
	CodeErr
	message string
}

// Error implements error interface for CodeErrWithContext
func (c *CodeErrWithContext) Error() string {
	return c.message
}

// ToGRPCStatus converts CodeErrWithContext to gRPC status
func (c *CodeErrWithContext) ToGRPCStatus() error {
	return status.Error(c.CodeErr.GetCodeErrEntity().GrpcCode, c.message)
}

// Unwrap returns the underlying CodeErr for errors.Is/As compatibility
func (c *CodeErrWithContext) Unwrap() error {
	return c.CodeErr
}

// IsErrorCode checks if an error is of a specific code, handling wrapped errors
func IsErrorCode(err error, code CodeErr) bool {
	if err == code {
		return true
	}
	// Handle CodeErrWithContext
	if contextErr, ok := err.(*CodeErrWithContext); ok {
		return contextErr.CodeErr == code
	}
	// Handle other wrapped errors using errors.Is
	var codeErr CodeErr
	if errors.As(err, &codeErr) {
		return codeErr == code
	}
	return false
}
