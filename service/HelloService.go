package service

import (
	"context"
	"fmt"

	hellopb "github.com/harryosmar/protobuf-go/gen/hello"
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
	message := fmt.Sprintf("Hello, %s!", req.GetName())
	return &hellopb.HelloResponse{
		Message: message,
	}, nil
}
