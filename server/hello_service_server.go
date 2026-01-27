package server

import (
	"context"

	appError "github.com/harryosmar/protobuf-go/error"
	hellopb "github.com/harryosmar/protobuf-go/gen/hello"
	"github.com/harryosmar/protobuf-go/logger"
	"github.com/harryosmar/protobuf-go/usecase"
	"go.uber.org/zap"
)

// HelloServiceServer implements the HelloService with usecase pattern
type HelloServiceServer struct {
	hellopb.UnimplementedHelloServiceServer
	helloUsecase usecase.HelloServiceUsecase
}

// NewHelloServiceServer creates a new HelloServiceServer instance
func NewHelloServiceServer() *HelloServiceServer {
	return &HelloServiceServer{
		helloUsecase: usecase.NewHelloServiceUsecase(),
	}
}

// GetHello implements the GetHello RPC method
func (s *HelloServiceServer) GetHello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	// Get logger with request ID from context
	log := logger.FromContext(ctx)
	log.Info("HelloService.GetHello called", zap.String("name", req.Name))

	// Validation will be handled by protoc-gen-validate generated code
	if err := req.Validate(); err != nil {
		return nil, appError.ErrInvalidArgument.WithMessage("validation failed: %v", err)
	}

	// Call usecase to handle business logic
	response, err := s.helloUsecase.GetHello(ctx, req)
	if err != nil {
		log.Error("Failed to process hello request", zap.String("name", req.Name), zap.Error(err))
		// Error conversion handled automatically by ErrorConversionInterceptor
		return nil, err
	}

	log.Info("HelloService.GetHello response", zap.String("message", response.Message))
	return response, nil
}
