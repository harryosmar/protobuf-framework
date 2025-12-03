package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/harryosmar/protobuf-go/gen"
	"github.com/harryosmar/protobuf-go/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcPort = ":50051"
	httpPort = ":8080"
)

func main() {
	// Start gRPC server in a goroutine
	go func() {
		if err := runGRPCServer(); err != nil {
			log.Fatalf("Failed to run gRPC server: %v", err)
		}
	}()

	// Start HTTP gateway server
	if err := runHTTPGateway(); err != nil {
		log.Fatalf("Failed to run HTTP gateway: %v", err)
	}
}

func runGRPCServer() error {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterHelloServiceServer(grpcServer, server.NewHelloServer())

	log.Printf("gRPC server listening on %s", grpcPort)
	return grpcServer.Serve(lis)
}

func runHTTPGateway() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := pb.RegisterHelloServiceHandlerFromEndpoint(ctx, mux, "localhost"+grpcPort, opts)
	if err != nil {
		return err
	}

	log.Printf("HTTP gateway listening on %s", httpPort)
	return http.ListenAndServe(httpPort, mux)
}
