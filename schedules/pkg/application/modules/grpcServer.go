package modules

import (
	"context"
	"errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type GrpcServer struct {
}

func (s *GrpcServer) Run(
	ctx context.Context,
	g *errgroup.Group,
	grpcServer *grpc.Server,
	listener net.Listener,
) {
	g.Go(func() error {
		go func() {
			<-ctx.Done()
			grpcServer.GracefulStop()
			logger(ctx).Info("gRPC server stopped gracefully")
		}()
		logger(ctx).Info("grpc server started", slog.String("address", listener.Addr().String()))
		if err := grpcServer.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			return err
		}
		return nil
	})
}
