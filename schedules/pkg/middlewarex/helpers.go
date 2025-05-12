package middlewarex

import (
	"context"
	"github.com/tummerk/golang/schedules/pkg/contextx"
	"log/slog"
	"net/http"
	"time"
)

func RequestLogRestapi(r *http.Request) slog.Attr {
	query := r.URL.Query()
	maskUserID, _ := contextx.MaskUserIDFromContext(r.Context())
	scheduleID := query.Get("schedule_id")
	requestInfo := []slog.Attr{
		slog.Any("traceID", contextx.TraceIDFromContext(r.Context())),
		slog.String("method", r.Method),
		slog.String("path", r.URL.String()),
		slog.String("host", r.Host),
		slog.String("user_agent", r.UserAgent()),
		slog.String("ip", r.RemoteAddr),
		slog.String("userID", maskUserID.String()),
	}
	if scheduleID != "" {
		requestInfo = append(requestInfo, slog.String("scheduleID", scheduleID))
	}
	return slog.Any("request_info", requestInfo)
}

func ResponseLogRestapi(ctx context.Context, w LoggingResponseWriter) slog.Attr {
	start := ctx.Value("startTime")
	startTime := start.(time.Time)
	responseInfo := []slog.Attr{
		slog.Any("traceID", contextx.TraceIDFromContext(ctx)),
		slog.Int("status", w.StatusCode),
		slog.Int("Size", w.Size),
		slog.Int64("duration", time.Since(startTime).Milliseconds()),
	}
	return slog.Any("response_info", responseInfo)
}

type LoggingResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	Size       int
}

func (lw *LoggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lw.ResponseWriter.Write(b)
	if err != nil {
		return size, err
	}
	lw.Size += size
	return size, nil
}

func (lw *LoggingResponseWriter) WriteHeader(statusCode int) {
	lw.ResponseWriter.WriteHeader(statusCode)
	lw.StatusCode = statusCode
}
