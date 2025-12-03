package server

import (
	"context"
	"fmt"

	pb "github.com/harryosmar/protobuf-go/gen"
)

// HelloServer implements the HelloService
type HelloServer struct {
	pb.UnimplementedHelloServiceServer
}

// NewHelloServer creates a new HelloServer instance
func NewHelloServer() *HelloServer {
	return &HelloServer{}
}

// GetHello implements the GetHello RPC method
func (s *HelloServer) GetHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	message := fmt.Sprintf("Hello, %s!", req.GetName())
	return &pb.HelloResponse{
		Message: message,
	}, nil
}
