package middleware

import (
	"context"
	"encoding/json"
	"fmt"
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

		// Conditionally serialize request payload for high-traffic optimization
		var reqPayload []byte
		var reqFields []zap.Field

		// Only log payload for non-high-frequency methods or small payloads
		if shouldLogPayload(info.FullMethod, req) {
			reqPayload, _ = json.Marshal(req)
			reqFields = []zap.Field{
				zap.String("method", info.FullMethod),
				zap.ByteString("request_payload", reqPayload),
				zap.Time("start_time", startTime),
			}
		} else {
			reqFields = []zap.Field{
				zap.String("method", info.FullMethod),
				zap.String("payload_size", getPayloadSize(req)),
				zap.Time("start_time", startTime),
			}
		}

		// Log request start
		log.Info("gRPC request received", reqFields...)

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

// shouldLogPayload determines if we should log the full payload based on method and size
func shouldLogPayload(method string, payload interface{}) bool {
	// Skip payload logging for high-frequency methods
	highFrequencyMethods := []string{
		"/grpc.health.v1.Health/Check",
		"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo",
	}

	for _, hfMethod := range highFrequencyMethods {
		if method == hfMethod {
			return false
		}
	}

	// Check payload size (rough estimation)
	payloadStr := getPayloadSize(payload)
	if len(payloadStr) > 1000 { // If estimated size > 1KB, skip full payload
		return false
	}

	return true
}

// getPayloadSize returns a string representation of payload size
func getPayloadSize(payload interface{}) string {
	if payload == nil {
		return "0B"
	}

	// Quick size estimation without full marshaling
	switch v := payload.(type) {
	case string:
		return fmt.Sprintf("%dB", len(v))
	default:
		// For other types, use a rough estimation
		return "~unknown"
	}
}
