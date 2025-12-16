package service

import (
	"context"
	"fmt"

	hellopb "github.com/harryosmar/protobuf-go/gen/hello"
	"github.com/harryosmar/protobuf-go/logger"
	"go.uber.org/zap"
)

// HelloServiceServer implements the HelloService
type HelloServiceServer struct {
	hellopb.UnimplementedHelloServiceServer
}

// NewHelloServiceServer creates a new HelloServiceServer instance
func NewHelloServiceServer() *HelloServiceServer {
	return &HelloServiceServer{}
}

// GetHello implements the GetHello RPC method
func (s *HelloServiceServer) GetHello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	// Get logger with request ID from context
	log := logger.FromContext(ctx)
	log.Info("HelloService.GetHello called", zap.String("name", req.GetName()))

	message := fmt.Sprintf("Hello, %s!", req.GetName())
	return &hellopb.HelloResponse{
		Message: message,
	}, nil
}
