package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	// Application settings
	AppName    string `envconfig:"APP_NAME" default:"protobuf-go-server"`
	AppVersion string `envconfig:"APP_VERSION" default:"v1.0.0"`

	// Server ports
	GRPCPort string `envconfig:"GRPC_PORT" default:":50051"`
	HTTPPort string `envconfig:"HTTP_PORT" default:":8080"`

	// Rate limiting configuration
	RateLimitEnabled        bool   `envconfig:"RATE_LIMIT_ENABLED" default:"true"`
	RateLimitRequestsPerSec int    `envconfig:"RATE_LIMIT_REQUESTS_PER_SEC" default:"100"`
	RateLimitBurstSize      int    `envconfig:"RATE_LIMIT_BURST_SIZE" default:"200"`
	RateLimitStrategy       string `envconfig:"RATE_LIMIT_STRATEGY" default:"global"` // global, per-method
}

// Get loads configuration from environment variables
func Get() *Config {
	cfg := Config{}
	envconfig.MustProcess("", &cfg)

	return &cfg
}

// GetRateLimitConfig returns rate limiting configuration
func (c *Config) GetRateLimitConfig() (requestsPerSec, burstSize int, strategy string) {
	return c.RateLimitRequestsPerSec, c.RateLimitBurstSize, c.RateLimitStrategy
}
