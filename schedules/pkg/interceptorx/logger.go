package interceptorx

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/tummerk/golang/schedules/pkg/contextx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"log/slog"
	"time"
)

func LoggingInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()

		traceID := slog.Any("traceID", contextx.TraceIDFromContext(ctx))

		logger.Info("request", traceID, slog.String("method", info.FullMethod))
		res, err := handler(ctx, req)
		if err != nil {
			logger.Error(err.Error(), traceID, slog.String("method", info.FullMethod))
		}
		duration := time.Since(start).Milliseconds()
		st, _ := status.FromError(err)
		statusCode := st.Code()
		size := proto.Size(res.(proto.Message))
		logger.Info("response", traceID, slog.Any("response_info",
			[]slog.Attr{
				slog.String("method", info.FullMethod),
				slog.Int("duration", int(duration)),
				slog.Int("size", size),
				slog.String("status", statusCode.String()),
			}))
		return res, err
	}
}
