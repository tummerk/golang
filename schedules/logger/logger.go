package logger

import (
	"context"
	"github.com/tummerk/golang/schedules/config"
	"github.com/tummerk/golang/schedules/utils"
	"gopkg.in/natefinch/lumberjack.v2"
	"log/slog"
	"net/http"
	"time"
)

func NewLogger(fileName string, level slog.Level) *slog.Logger {
	logRotator := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     30,
	}
	defer logRotator.Close()

	handler := slog.NewJSONHandler(logRotator, &slog.HandlerOptions{
		Level: level,
	})
	logger := slog.New(handler)
	return logger
}

func RequestLog(r *http.Request) slog.Attr {
	query := r.URL.Query()
	userID := query.Get("user_id")
	scheduleID := query.Get("schedule_id")
	userID, _ = utils.Encrypt(userID, config.Key)
	requestInfo := []slog.Attr{
		slog.Any("traceID", r.Context().Value("traceID")),
		slog.String("method", r.Method),
		slog.String("path", r.URL.String()),
		slog.String("host", r.Host),
		slog.String("user_agent", r.UserAgent()),
		slog.String("ip", r.RemoteAddr),
		slog.String("userID", userID),
	}
	if scheduleID != "" {
		requestInfo = append(requestInfo, slog.String("scheduleID", scheduleID))
	}
	return slog.Any("request_info", requestInfo)
}

func ResponseLog(ctx context.Context, w LoggingResponseWriter) slog.Attr {
	start := ctx.Value("startTime")
	startTime := start.(time.Time)
	responseInfo := []slog.Attr{
		slog.Any("traceID", ctx.Value("traceID")),
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
