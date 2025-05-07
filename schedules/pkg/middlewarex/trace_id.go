package middlewarex

import (
	"github.com/tummerk/golang/schedules/pkg/contextx"
	"net/http"

	"github.com/rs/xid"
)

const headerNameTraceID = "X-Trace-Id"

func TraceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get(headerNameTraceID)

		if traceID == "" {
			traceID = xid.New().String()
		}

		ctx := contextx.WithTraceID(r.Context(), contextx.TraceID(traceID))

		w.Header().Set(headerNameTraceID, traceID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
