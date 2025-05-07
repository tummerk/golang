package interceptorx

import (
	"context"
	"github.com/rs/xid"
	"github.com/tummerk/golang/schedules/pkg/contextx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TraceIDInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		traceID := ""
		if traceIDs := md.Get("x-trace-id"); len(traceIDs) > 0 {
			traceID = traceIDs[0]
		} else {
			traceID = xid.New().String()
		}

		newCtx := contextx.WithTraceID(ctx, contextx.TraceID(traceID))
		return handler(newCtx, req)
	}
}
