package middleware

import (
	"context"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// gRPC request metrics
	grpcRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"method", "status_code"},
	)

	grpcRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_request_duration_seconds",
			Help:    "Duration of gRPC requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status_code"},
	)

	grpcActiveConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "grpc_active_connections",
			Help: "Number of active gRPC connections",
		},
	)

	// Rate limiting metrics
	rateLimitExceeded = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rate_limit_exceeded_total",
			Help: "Total number of rate limit exceeded events",
		},
		[]string{"method", "key"},
	)
)

// MetricsInterceptor collects Prometheus metrics for gRPC requests
func MetricsInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()

		// Increment active connections
		grpcActiveConnections.Inc()
		defer grpcActiveConnections.Dec()

		// Call the handler
		resp, err := handler(ctx, req)

		// Calculate duration
		duration := time.Since(startTime)

		// Determine gRPC status code
		statusCode := codes.OK
		if err != nil {
			if st, ok := status.FromError(err); ok {
				statusCode = st.Code()
			} else {
				statusCode = codes.Internal
			}
		}

		// Record metrics
		labels := prometheus.Labels{
			"method":      info.FullMethod,
			"status_code": strconv.Itoa(int(statusCode)),
		}

		grpcRequestsTotal.With(labels).Inc()
		grpcRequestDuration.With(labels).Observe(duration.Seconds())

		return resp, err
	}
}

// RecordRateLimitExceeded records rate limit exceeded events
func RecordRateLimitExceeded(method, key string) {
	rateLimitExceeded.With(prometheus.Labels{
		"method": method,
		"key":    key,
	}).Inc()
}
