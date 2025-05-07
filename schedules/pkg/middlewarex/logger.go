package middlewarex

import (
	"context"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		logger(r.Context()).Info("request", RequestLogRestapi(r))

		ctx := context.WithValue(r.Context(), "startTime", startTime)
		ctx = context.WithValue(ctx, "ip", r.RemoteAddr)
		ctx = context.WithValue(ctx, "userAgent", r.UserAgent())
		r = r.WithContext(ctx)

		lw := LoggingResponseWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
			Size:           0,
		}
		next.ServeHTTP(&lw, r)
		logger(r.Context()).Info("response", ResponseLogRestapi(r.Context(), lw))
	})
}
