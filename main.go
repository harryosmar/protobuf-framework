package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	hellopb "github.com/harryosmar/protobuf-go/gen/hello"
	userpb "github.com/harryosmar/protobuf-go/gen/user"
	"github.com/harryosmar/protobuf-go/logger"
	"github.com/harryosmar/protobuf-go/middleware"
	"github.com/harryosmar/protobuf-go/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcPort = ":50051"
	httpPort = ":8080"
)

func main() {
	// Initialize logger
	baseLogger, err := logger.InitLogger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer baseLogger.Sync()

	// Start gRPC server in a goroutine
	go func() {
		if err := runGRPCServer(baseLogger); err != nil {
			baseLogger.Fatal("Failed to run gRPC server", zap.Error(err))
		}
	}()

	// Start HTTP gateway server
	if err := runHTTPGateway(baseLogger); err != nil {
		baseLogger.Fatal("Failed to run HTTP gateway", zap.Error(err))
	}
}

func runGRPCServer(baseLogger *zap.Logger) error {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RequestIDInterceptor(baseLogger),
			middleware.NewGlobalRateLimitInterceptor(100, 200), // 100 req/sec, 200 burst
			middleware.LoggingInterceptor(baseLogger),
			// Add future interceptors here (auth, metrics, etc.)
		),
	)
	hellopb.RegisterHelloServiceServer(grpcServer, service.NewHelloServer())
	userpb.RegisterUserServiceServer(grpcServer, service.NewUserServer())

	baseLogger.Info("gRPC server listening", zap.String("port", grpcPort))
	return grpcServer.Serve(lis)
}

func runHTTPGateway(baseLogger *zap.Logger) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := hellopb.RegisterHelloServiceHandlerFromEndpoint(ctx, mux, "localhost"+grpcPort, opts)
	if err != nil {
		return err
	}

	err = userpb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "localhost"+grpcPort, opts)
	if err != nil {
		return err
	}

	baseLogger.Info("HTTP gateway listening", zap.String("port", httpPort))
	return http.ListenAndServe(httpPort, mux)
}
