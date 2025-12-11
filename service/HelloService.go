package service

import (
	"context"
	"fmt"

	hellopb "github.com/harryosmar/protobuf-go/gen/hello"
	"github.com/harryosmar/protobuf-go/logger"
	"go.uber.org/zap"
)

// HelloServer implements the HelloService
type HelloServer struct {
	hellopb.UnimplementedHelloServiceServer
}

// NewHelloServer creates a new HelloServer instance
func NewHelloServer() *HelloServer {
	return &HelloServer{}
}

// GetHello implements the GetHello RPC method
func (s *HelloServer) GetHello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	// Get logger with request ID from context
	log := logger.FromContext(ctx)
	log.Info("HelloService.GetHello called", zap.String("name", req.GetName()))

	message := fmt.Sprintf("Hello, %s!", req.GetName())
	return &hellopb.HelloResponse{
		Message: message,
	}, nil
}
