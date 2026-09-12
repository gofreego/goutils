package metrics

import (
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_server_requests_total",
			Help: "Total HTTP requests handled, labeled by method, normalized path, and status code.",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_server_request_duration_seconds",
			Help:    "HTTP request handling duration in seconds, labeled by method and normalized path.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

var (
	numericSegment = regexp.MustCompile(`^[0-9]+$`)
	uuidSegment    = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
)

// normalizePath collapses path segments that look like IDs (numeric or
// UUID) into a fixed placeholder, so /exams/42/attempt and
// /exams/1337/attempt count as the same series instead of each request's
// ID value creating its own metric series forever. This is a heuristic,
// not a real router-aware template (grpc-gateway doesn't expose its
// matched route pattern to wrapping middleware) - it errs on the side of
// collapsing anything ID-shaped rather than trying to be exact.
func normalizePath(path string) string {
	segments := make([]byte, 0, len(path))
	start := 0
	for i := 0; i <= len(path); i++ {
		if i < len(path) && path[i] != '/' {
			continue
		}
		seg := path[start:i]
		if seg != "" {
			if numericSegment.MatchString(seg) || uuidSegment.MatchString(seg) {
				seg = ":id"
			}
		}
		segments = append(segments, '/')
		segments = append(segments, seg...)
		start = i + 1
	}
	if len(segments) == 0 {
		return "/"
	}
	return string(segments[1:])
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// WithHTTPMetrics wraps a handler to record request count and latency for
// every HTTP request, labeled by method, normalized path, and status code.
// Use this for the REST gateway - the gRPC UnaryServerInterceptor in this
// package does NOT see this traffic, since RegisterXxxHandlerServer calls
// service methods in-process rather than through the real grpc.Server.
func WithHTTPMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		path := normalizePath(r.URL.Path)
		status := strconv.Itoa(rec.status)

		httpRequestsTotal.WithLabelValues(r.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(r.Method, path).Observe(time.Since(start).Seconds())
	})
}
