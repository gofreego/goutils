// Package metrics provides Prometheus instrumentation shared across
// services. Request metrics are recorded at the gRPC layer rather than the
// HTTP layer: every service here exposes gRPC and forwards its REST
// gateway into the same in-process gRPC methods, so one interceptor covers
// both, and the gRPC method name (e.g.
// /learnerservice_v1.LearnerService/GetTopics) is a natural low-cardinality
// route label - unlike the raw HTTP path, it never varies with path
// parameters, so it can't cause a cardinality blowup the way labeling by
// resolved URL would.
package metrics

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	requestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_server_requests_total",
			Help: "Total gRPC requests handled, labeled by method and status code.",
		},
		[]string{"method", "code"},
	)

	requestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_server_request_duration_seconds",
			Help:    "gRPC request handling duration in seconds, labeled by method.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
)

// UnaryServerInterceptor records a count and duration for every unary gRPC
// call. Wire it in with grpc.NewServer(grpc.UnaryInterceptor(metrics.UnaryServerInterceptor())).
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)

		code := status.Code(err).String()
		requestsTotal.WithLabelValues(info.FullMethod, code).Inc()
		requestDuration.WithLabelValues(info.FullMethod).Observe(time.Since(start).Seconds())

		return resp, err
	}
}

// Handler serves the Prometheus text exposition format for scraping.
func Handler() http.Handler {
	return promhttp.Handler()
}
