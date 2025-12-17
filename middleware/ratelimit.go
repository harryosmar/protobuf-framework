package middleware

import (
	"context"
	error2 "github.com/harryosmar/protobuf-go/error"
	"sync"

	"github.com/harryosmar/protobuf-go/config"
	"github.com/harryosmar/protobuf-go/logger"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	RequestsPerSecond int          // Number of requests allowed per second
	BurstSize         int          // Maximum burst size
	KeyExtractor      KeyExtractor // Function to extract rate limit key from context
}

// KeyExtractor extracts a key from context for rate limiting (e.g., client IP, user ID)
type KeyExtractor func(ctx context.Context, info *grpc.UnaryServerInfo) string

// DefaultKeyExtractor uses the method name as the rate limit key (global rate limit)
func DefaultKeyExtractor(ctx context.Context, info *grpc.UnaryServerInfo) string {
	return "global"
}

// MethodKeyExtractor uses the gRPC method as the rate limit key (per-method rate limiting)
func MethodKeyExtractor(ctx context.Context, info *grpc.UnaryServerInfo) string {
	return info.FullMethod
}

// RateLimiter manages rate limiters for different keys
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	config   RateLimitConfig
	mutex    sync.RWMutex
}

// NewRateLimiter creates a new rate limiter instance
func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	if config.KeyExtractor == nil {
		config.KeyExtractor = DefaultKeyExtractor
	}
	if config.RequestsPerSecond <= 0 {
		config.RequestsPerSecond = 100 // Default: 100 requests per second
	}
	if config.BurstSize <= 0 {
		config.BurstSize = config.RequestsPerSecond * 2 // Default: 2x burst
	}

	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		config:   config,
	}
}

// getLimiter gets or creates a rate limiter for the given key
func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mutex.RLock()
	limiter, exists := rl.limiters[key]
	rl.mutex.RUnlock()

	if exists {
		return limiter
	}

	// Create new limiter if it doesn't exist
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	// Double-check after acquiring write lock
	if limiter, exists := rl.limiters[key]; exists {
		return limiter
	}

	// Create new rate limiter
	limiter = rate.NewLimiter(rate.Limit(rl.config.RequestsPerSecond), rl.config.BurstSize)
	rl.limiters[key] = limiter
	return limiter
}

// RateLimitInterceptor creates a gRPC interceptor for rate limiting
func RateLimitInterceptor(rateLimiter *RateLimiter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Extract rate limit key
		key := rateLimiter.config.KeyExtractor(ctx, info)

		// Get rate limiter for this key
		limiter := rateLimiter.getLimiter(key)

		// Check if request is allowed
		if !limiter.Allow() {
			// Get logger from context for rate limit logging
			log := logger.FromContext(ctx)
			log.Warn("Rate limit exceeded",
				zap.String("method", info.FullMethod),
				zap.String("rate_limit_key", key),
				zap.Int("requests_per_second", rateLimiter.config.RequestsPerSecond),
				zap.Int("burst_size", rateLimiter.config.BurstSize),
			)

			// Record rate limit exceeded metric
			RecordRateLimitExceeded(info.FullMethod, key)

			// Return rate limit exceeded error
			return nil, error2.ErrResourceExhausted.WithMessage(
				"Rate limit exceeded. Maximum %d requests per second allowed.",
				rateLimiter.config.RequestsPerSecond)
		}

		// Request allowed, proceed to handler
		return handler(ctx, req)
	}
}

// NewGlobalRateLimitInterceptor creates a rate limiter with global limits
func NewGlobalRateLimitInterceptor(requestsPerSecond, burstSize int) grpc.UnaryServerInterceptor {
	config := RateLimitConfig{
		RequestsPerSecond: requestsPerSecond,
		BurstSize:         burstSize,
		KeyExtractor:      DefaultKeyExtractor,
	}
	rateLimiter := NewRateLimiter(config)
	return RateLimitInterceptor(rateLimiter)
}

// NewPerMethodRateLimitInterceptor creates a rate limiter with per-method limits
func NewPerMethodRateLimitInterceptor(requestsPerSecond, burstSize int) grpc.UnaryServerInterceptor {
	config := RateLimitConfig{
		RequestsPerSecond: requestsPerSecond,
		BurstSize:         burstSize,
		KeyExtractor:      MethodKeyExtractor,
	}
	rateLimiter := NewRateLimiter(config)
	return RateLimitInterceptor(rateLimiter)
}

func NewRateLimitInterceptors(cfg *config.Config) []grpc.UnaryServerInterceptor {
	if !cfg.RateLimitEnabled {
		return []grpc.UnaryServerInterceptor{}
	}

	if cfg.RateLimitStrategy == "per-method" {
		return []grpc.UnaryServerInterceptor{
			NewPerMethodRateLimitInterceptor(cfg.RateLimitRequestsPerSec, cfg.RateLimitBurstSize),
		}
	}

	return []grpc.UnaryServerInterceptor{
		NewGlobalRateLimitInterceptor(cfg.RateLimitRequestsPerSec, cfg.RateLimitBurstSize),
	}
}
