package usecase

import (
	"context"
	"fmt"

	hellopb "github.com/harryosmar/protobuf-go/gen/hello"
)

// HelloServiceUsecase defines the interface for hello business logic
type HelloServiceUsecase interface {
	GetHello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error)
}

// helloServiceUsecase implements HelloServiceUsecase interface
type helloServiceUsecase struct{}

// NewHelloServiceUsecase creates a new hello usecase instance
func NewHelloServiceUsecase() HelloServiceUsecase {
	return &helloServiceUsecase{}
}

// GetHello handles the business logic for greeting a user
func (u *helloServiceUsecase) GetHello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	// Simple business logic to create a greeting message
	message := fmt.Sprintf("Hello, %s!", req.Name)
	
	return &hellopb.HelloResponse{
		Message: message,
	}, nil
}
