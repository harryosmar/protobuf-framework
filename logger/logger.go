package logger

import (
	"context"

	"go.uber.org/zap"
)

const LoggerContextKey = "logger"

// InitLogger creates a new zap logger instance
func InitLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	return config.Build()
}

// WithRequestID creates a logger with request_id field pre-populated
func WithRequestID(logger *zap.Logger, requestID string) *zap.Logger {
	return logger.With(zap.String("request_id", requestID))
}

// FromContext extracts the logger from context
func FromContext(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(LoggerContextKey).(*zap.Logger); ok {
		return logger
	}
	// Fallback to a basic logger if not found in context
	logger, _ := zap.NewProduction()
	return logger
}

// ToContext adds logger to context
func ToContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, LoggerContextKey, logger)
}
