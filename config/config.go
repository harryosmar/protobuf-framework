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

	// Database configuration
	DatabaseURL     string `envconfig:"DATABASE_URL" default:"root:password@tcp(localhost:3306)/protobuf_go?charset=utf8mb4&parseTime=True&loc=Local"`
	DatabaseMaxIdle int    `envconfig:"DATABASE_MAX_IDLE" default:"10"`
	DatabaseMaxOpen int    `envconfig:"DATABASE_MAX_OPEN" default:"100"`
	DatabaseMaxLife int    `envconfig:"DATABASE_MAX_LIFE" default:"3600"` // seconds

	// Rate limiting configuration
	RateLimitEnabled        bool   `envconfig:"RATE_LIMIT_ENABLED" default:"true"`
	RateLimitRequestsPerSec int    `envconfig:"RATE_LIMIT_REQUESTS_PER_SEC" default:"100"`
	RateLimitBurstSize      int    `envconfig:"RATE_LIMIT_BURST_SIZE" default:"200"`
	RateLimitStrategy       string `envconfig:"RATE_LIMIT_STRATEGY" default:"global"` // global, per-method

	// gRPC server configuration
	GRPCMaxConnectionIdle     int  `envconfig:"GRPC_MAX_CONNECTION_IDLE" default:"15"`     // seconds
	GRPCMaxConnectionAge      int  `envconfig:"GRPC_MAX_CONNECTION_AGE" default:"30"`      // seconds
	GRPCMaxConnectionAgeGrace int  `envconfig:"GRPC_MAX_CONNECTION_AGE_GRACE" default:"5"` // seconds
	GRPCKeepaliveTime         int  `envconfig:"GRPC_KEEPALIVE_TIME" default:"5"`           // seconds
	GRPCKeepaliveTimeout      int  `envconfig:"GRPC_KEEPALIVE_TIMEOUT" default:"1"`        // seconds
	GRPCKeepaliveMinTime      int  `envconfig:"GRPC_KEEPALIVE_MIN_TIME" default:"5"`       // seconds
	GRPCPermitWithoutStream   bool `envconfig:"GRPC_PERMIT_WITHOUT_STREAM" default:"false"`
	GRPCMaxRecvMsgSize        int  `envconfig:"GRPC_MAX_RECV_MSG_SIZE" default:"4194304"` // 4MB in bytes
	GRPCMaxSendMsgSize        int  `envconfig:"GRPC_MAX_SEND_MSG_SIZE" default:"4194304"` // 4MB in bytes
	GRPCMaxConcurrentStreams  int  `envconfig:"GRPC_MAX_CONCURRENT_STREAMS" default:"1000"`
}

// Get loads configuration from environment variables
func Get() *Config {
	cfg := Config{}
	envconfig.MustProcess("", &cfg)

	return &cfg
}
