package middleware

import (
	"log/slog"
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

func logRequestError(r *http.Request, message string, statusCode int, err error, attrs ...slog.Attr) {
	baseAttrs := []slog.Attr{
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.Int("status", statusCode),
	}

	if err != nil {
		baseAttrs = append(baseAttrs, slog.Any("error", err))
	}

	if spanContext := trace.SpanContextFromContext(r.Context()); spanContext.HasTraceID() {
		baseAttrs = append(baseAttrs,
			slog.String("trace_id", spanContext.TraceID().String()),
			slog.String("span_id", spanContext.SpanID().String()),
		)
	}

	baseAttrs = append(baseAttrs, attrs...)
	slog.LogAttrs(r.Context(), logLevelForStatus(statusCode), message, baseAttrs...)
}

func logLevelForStatus(statusCode int) slog.Level {
	if statusCode >= http.StatusInternalServerError {
		return slog.LevelError
	}
	return slog.LevelWarn
}
