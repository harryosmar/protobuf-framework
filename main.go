package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/harryosmar/protobuf-go/config"
	hellopb "github.com/harryosmar/protobuf-go/gen/hello"
	userpb "github.com/harryosmar/protobuf-go/gen/user"
	"github.com/harryosmar/protobuf-go/logger"
	"github.com/harryosmar/protobuf-go/middleware"
	"github.com/harryosmar/protobuf-go/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Load configuration
	cfg := config.Get()

	// Initialize logger
	baseLogger, err := logger.InitLogger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer baseLogger.Sync()

	baseLogger.Info("Starting server",
		zap.String("app_name", cfg.AppName),
		zap.String("app_version", cfg.AppVersion),
		zap.String("grpc_port", cfg.GRPCPort),
		zap.String("http_port", cfg.HTTPPort),
	)

	// Start gRPC server in a goroutine
	go func() {
		if err := runGRPCServer(cfg, baseLogger); err != nil {
			baseLogger.Fatal("Failed to run gRPC server", zap.Error(err))
		}
	}()

	// Start HTTP gateway server
	if err := runHTTPGateway(cfg, baseLogger); err != nil {
		baseLogger.Fatal("Failed to run HTTP gateway", zap.Error(err))
	}
}

func runGRPCServer(cfg *config.Config, baseLogger *zap.Logger) error {
	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		return err
	}

	// Get rate limiting configuration
	reqPerSec, burstSize, strategy := cfg.GetRateLimitConfig()

	var rateLimitInterceptor grpc.UnaryServerInterceptor
	if cfg.RateLimitEnabled {
		if strategy == "per-method" {
			rateLimitInterceptor = middleware.NewPerMethodRateLimitInterceptor(reqPerSec, burstSize)
		} else {
			rateLimitInterceptor = middleware.NewGlobalRateLimitInterceptor(reqPerSec, burstSize)
		}
	}

	// Build interceptor chain
	interceptors := []grpc.UnaryServerInterceptor{
		middleware.RequestIDInterceptor(baseLogger),
	}

	if cfg.RateLimitEnabled {
		interceptors = append(interceptors, rateLimitInterceptor)
	}

	interceptors = append(interceptors, middleware.LoggingInterceptor(baseLogger))

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptors...),
	)
	hellopb.RegisterHelloServiceServer(grpcServer, service.NewHelloServer())
	userpb.RegisterUserServiceServer(grpcServer, service.NewUserServer())

	baseLogger.Info("gRPC server listening", zap.String("port", cfg.GRPCPort))
	return grpcServer.Serve(lis)
}

func runHTTPGateway(cfg *config.Config, baseLogger *zap.Logger) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := hellopb.RegisterHelloServiceHandlerFromEndpoint(ctx, mux, "localhost"+cfg.GRPCPort, opts)
	if err != nil {
		return err
	}

	err = userpb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "localhost"+cfg.GRPCPort, opts)
	if err != nil {
		return err
	}

	baseLogger.Info("HTTP gateway listening", zap.String("port", cfg.HTTPPort))
	return http.ListenAndServe(cfg.HTTPPort, mux)
}
