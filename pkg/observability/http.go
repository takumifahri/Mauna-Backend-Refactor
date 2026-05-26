package observability

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "mauna",
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "Total number of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "mauna",
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "HTTP request duration in seconds.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration)
}

// HTTPMiddleware adds request tracing and structured request logs.
func HTTPMiddleware(next http.Handler, logger *slog.Logger) http.Handler {
	return otelhttp.NewHandler(
		requestLogger(next, logger),
		"http.server",
		otelhttp.WithSpanNameFormatter(func(_ string, r *http.Request) string {
			return r.Method + " " + r.URL.Path
		}),
	)
}

func requestLogger(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		metrics := httpsnoop.CaptureMetrics(next, w, r)
		duration := time.Since(start)

		statusCode := metrics.Code
		if statusCode == 0 {
			statusCode = http.StatusOK
		}

		path := routePath(r)
		status := strconv.Itoa(statusCode)

		httpRequestsTotal.WithLabelValues(r.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(r.Method, path, status).Observe(duration.Seconds())

		spanContext := trace.SpanContextFromContext(r.Context())
		attrs := []slog.Attr{
			slog.String("method", r.Method),
			slog.String("path", path),
			slog.Int("status", statusCode),
			slog.Int64("bytes", metrics.Written),
			slog.Duration("duration", duration),
			slog.String("remote_addr", r.RemoteAddr),
		}

		if spanContext.HasTraceID() {
			attrs = append(attrs,
				slog.String("trace_id", spanContext.TraceID().String()),
				slog.String("span_id", spanContext.SpanID().String()),
			)
		}

		level := slog.LevelInfo
		if statusCode >= http.StatusInternalServerError {
			level = slog.LevelError
		} else if statusCode >= http.StatusBadRequest {
			level = slog.LevelWarn
		}

		logger.LogAttrs(r.Context(), level, "http_request", attrs...)
	})
}

func routePath(r *http.Request) string {
	if r.Pattern == "" {
		return r.URL.Path
	}

	if strings.HasPrefix(r.Pattern, r.Method+" ") {
		return strings.TrimPrefix(r.Pattern, r.Method+" ")
	}

	return r.Pattern
}
