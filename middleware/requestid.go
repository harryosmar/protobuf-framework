package middleware

import (
	"context"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const RequestIDHeader = "x-request-id"

// RequestIDInterceptor adds a request ID to gRPC requests if not present
func RequestIDInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Get metadata from context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}

	// Check if request ID already exists
	requestIDs := md.Get(RequestIDHeader)
	var requestID string

	if len(requestIDs) == 0 {
		// Generate new UUID if not present
		requestID = uuid.New().String()
		md.Set(RequestIDHeader, requestID)
		// Update context with new metadata
		ctx = metadata.NewIncomingContext(ctx, md)
	} else {
		requestID = requestIDs[0]
	}

	// Add request ID to outgoing metadata for response
	outgoingMD := metadata.Pairs(RequestIDHeader, requestID)
	ctx = metadata.NewOutgoingContext(ctx, outgoingMD)

	// Log the request with ID
	log.Printf("[%s] %s", requestID, info.FullMethod)

	// Call the handler
	resp, err := handler(ctx, req)

	// Set response header
	grpc.SetHeader(ctx, metadata.Pairs(RequestIDHeader, requestID))

	return resp, err
}
