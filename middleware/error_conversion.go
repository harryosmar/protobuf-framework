package middleware

import (
	"context"

	error2 "github.com/harryosmar/protobuf-go/error"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorConversionInterceptor automatically converts CodeErr to gRPC status
func ErrorConversionInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			// Convert CodeErr to gRPC status automatically
			if codeErr, ok := err.(error2.CodeErr); ok {
				return resp, codeErr.ToGRPCStatus()
			}
			if contextErr, ok := err.(*error2.CodeErrWithContext); ok {
				return resp, contextErr.ToGRPCStatus()
			}
			// For other errors, return as Internal error
			return resp, status.Error(codes.Internal, err.Error())
		}
		return resp, nil
	}
}
