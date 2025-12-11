package middleware

import (
	"context"

	"github.com/google/uuid"
	"github.com/harryosmar/protobuf-go/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	RequestIDHeader     = "x-request-id"
	RequestIDContextKey = "request-id"
)

// RequestIDInterceptor adds a request ID to gRPC requests if not present and injects logger with request ID
func RequestIDInterceptor(baseLogger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Get metadata from context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		// Check if request ID already exists
		requestIDs := md.Get(RequestIDHeader)
		var requestID string

		if len(requestIDs) == 0 {
			// Generate new UUID if not present
			requestID = uuid.New().String()
			md.Set(RequestIDHeader, requestID)
			// Update context with new metadata
			ctx = metadata.NewIncomingContext(ctx, md)
		} else {
			requestID = requestIDs[0]
		}

		// Store request ID in context for service access
		ctx = context.WithValue(ctx, RequestIDContextKey, requestID)

		// Create logger with request ID and add to context
		requestLogger := logger.WithRequestID(baseLogger, requestID)
		ctx = logger.ToContext(ctx, requestLogger)

		// Add request ID to outgoing metadata for response
		outgoingMD := metadata.Pairs(RequestIDHeader, requestID)
		ctx = metadata.NewOutgoingContext(ctx, outgoingMD)

		// Log the request with structured logging
		requestLogger.Info("gRPC request started", zap.String("method", info.FullMethod))

		// Call the handler
		resp, err := handler(ctx, req)

		// Set response header
		grpc.SetHeader(ctx, metadata.Pairs(RequestIDHeader, requestID))

		if err != nil {
			requestLogger.Error("gRPC request failed", zap.String("method", info.FullMethod), zap.Error(err))
		} else {
			requestLogger.Info("gRPC request completed", zap.String("method", info.FullMethod))
		}

		return resp, err
	}
}

// GetRequestID extracts the request ID from context
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDContextKey).(string); ok {
		return requestID
	}
	return ""
}
