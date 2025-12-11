package middleware

import (
	"context"
	"encoding/json"
	"time"

	"github.com/harryosmar/protobuf-go/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoggingInterceptor logs detailed request and response information
func LoggingInterceptor(baseLogger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()

		// Get logger from context (should have request_id from RequestIDInterceptor)
		log := logger.FromContext(ctx)

		// Serialize request payload
		reqPayload, _ := json.Marshal(req)

		// Log request start
		log.Info("gRPC request received",
			zap.String("method", info.FullMethod),
			zap.ByteString("request_payload", reqPayload),
			zap.Time("start_time", startTime),
		)

		// Call the handler
		resp, err := handler(ctx, req)

		// Calculate duration
		duration := time.Since(startTime)

		// Determine gRPC status
		grpcStatus := codes.OK
		statusCode := "OK"
		if err != nil {
			if st, ok := status.FromError(err); ok {
				grpcStatus = st.Code()
				statusCode = grpcStatus.String()
			} else {
				grpcStatus = codes.Internal
				statusCode = "INTERNAL"
			}
		}

		// Serialize response payload (if no error)
		var respPayload []byte
		if resp != nil {
			respPayload, _ = json.Marshal(resp)
		}

		// Common fields for logging
		fields := []zap.Field{
			zap.String("method", info.FullMethod),
			zap.String("grpc_status", statusCode),
			zap.Int("status_code", int(grpcStatus)),
			zap.Duration("duration", duration),
			zap.ByteString("request_payload", reqPayload),
		}

		// Add response payload if available
		if respPayload != nil {
			fields = append(fields, zap.ByteString("response_payload", respPayload))
		}

		// Log based on severity
		if err != nil {
			// Error severity for failed requests
			fields = append(fields, zap.Error(err))
			log.Error("gRPC request failed", fields...)
		} else {
			// Info severity for successful requests
			log.Info("gRPC request completed", fields...)
		}

		return resp, err
	}
}
