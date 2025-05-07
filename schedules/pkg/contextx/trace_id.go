package contextx

import (
	"context"
)

type TraceID string

type contextKeyTraceID struct{}

func (t TraceID) String() string {
	return string(t)
}

func WithTraceID(ctx context.Context, traceID TraceID) context.Context {
	return context.WithValue(ctx, contextKeyTraceID{}, traceID)
}

func TraceIDFromContext(ctx context.Context) TraceID {
	traceID, ok := ctx.Value(contextKeyTraceID{}).(TraceID)
	if !ok {
		return ""
	}

	return traceID
}
