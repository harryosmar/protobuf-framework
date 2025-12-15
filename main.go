package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/harryosmar/protobuf-go/config"
	"github.com/harryosmar/protobuf-go/database"
	hellopb "github.com/harryosmar/protobuf-go/gen/hello"
	userpb "github.com/harryosmar/protobuf-go/gen/user"
	"github.com/harryosmar/protobuf-go/handlers"
	"github.com/harryosmar/protobuf-go/logger"
	"github.com/harryosmar/protobuf-go/middleware"
	"github.com/harryosmar/protobuf-go/models"
	"github.com/harryosmar/protobuf-go/repository"
	"github.com/harryosmar/protobuf-go/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
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

	// Initialize database with new pattern
	db, err := database.NewDatabase(cfg, baseLogger)
	if err != nil {
		baseLogger.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer func() {
		if err := database.CloseDatabase(db); err != nil {
			baseLogger.Error("Failed to close database", zap.Error(err))
		}
	}()

	// Auto-migrate database schema
	if err := db.AutoMigrate(&models.User{}); err != nil {
		baseLogger.Fatal("Failed to migrate database", zap.Error(err))
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	baseLogger.Info("Starting server",
		zap.String("app_name", cfg.AppName),
		zap.String("app_version", cfg.AppVersion),
		zap.String("grpc_port", cfg.GRPCPort),
		zap.String("http_port", cfg.HTTPPort),
	)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Channel to listen for interrupt signal to trigger shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start gRPC server in a goroutine
	grpcDone := make(chan error, 1)
	go func() {
		grpcDone <- runGRPCServer(ctx, cfg, baseLogger, userRepo)
	}()

	// Start HTTP gateway server in a goroutine
	httpDone := make(chan error, 1)
	go func() {
		httpDone <- runHTTPGateway(cfg, baseLogger)
	}()

	// Wait for interrupt signal or server error
	select {
	case <-quit:
		baseLogger.Info("Shutdown signal received")
	case err := <-grpcDone:
		if err != nil {
			baseLogger.Error("gRPC server error", zap.Error(err))
		}
	case err := <-httpDone:
		if err != nil {
			baseLogger.Error("HTTP server error", zap.Error(err))
		}
	}

	// Graceful shutdown
	baseLogger.Info("Shutting down servers...")
	cancel()

	// Give servers time to shutdown gracefully
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	select {
	case <-shutdownCtx.Done():
		baseLogger.Warn("Shutdown timeout exceeded")
	case <-time.After(5 * time.Second):
		baseLogger.Info("Servers shutdown completed")
	}
}

func runGRPCServer(ctx context.Context, cfg *config.Config, baseLogger *zap.Logger, userRepo repository.UserRepository) error {
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
		middleware.MetricsInterceptor(), // Add metrics collection
	}

	if cfg.RateLimitEnabled {
		interceptors = append(interceptors, rateLimitInterceptor)
	}

	interceptors = append(interceptors, middleware.LoggingInterceptor(baseLogger))

	// Production-ready gRPC server with keepalive and limits
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptors...),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     15 * time.Second,
			MaxConnectionAge:      30 * time.Second,
			MaxConnectionAgeGrace: 5 * time.Second,
			Time:                  5 * time.Second,
			Timeout:               1 * time.Second,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second,
			PermitWithoutStream: false,
		}),
		grpc.MaxRecvMsgSize(4*1024*1024), // 4MB
		grpc.MaxSendMsgSize(4*1024*1024), // 4MB
		grpc.MaxConcurrentStreams(1000),
	)

	hellopb.RegisterHelloServiceServer(grpcServer, service.NewHelloServer())
	userpb.RegisterUserServiceServer(grpcServer, service.NewUserServer(userRepo))

	baseLogger.Info("gRPC server listening", zap.String("port", cfg.GRPCPort))

	// Graceful shutdown handling
	go func() {
		<-ctx.Done()
		baseLogger.Info("Gracefully stopping gRPC server...")
		grpcServer.GracefulStop()
	}()

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

	// Create a new HTTP mux for additional endpoints
	httpMux := http.NewServeMux()

	// Register gRPC gateway
	httpMux.Handle("/", mux)

	// Register health endpoint
	httpMux.HandleFunc("/health", handlers.HealthHandler(cfg))

	// Register Swagger endpoints
	httpMux.HandleFunc("/docs", handlers.SwaggerUIHandler())
	httpMux.HandleFunc("/docs/swagger.json", handlers.SwaggerHandler())

	// Register Prometheus metrics endpoint
	httpMux.Handle("/metrics", promhttp.Handler())

	baseLogger.Info("HTTP gateway listening", zap.String("port", cfg.HTTPPort))
	return http.ListenAndServe(cfg.HTTPPort, httpMux)
}
