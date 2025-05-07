package modules

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"net/http"
	"time"
)

type HTTPServer struct {
	ShutdownTimeout time.Duration
}

func (h HTTPServer) Run(
	ctx context.Context,
	g *errgroup.Group,
	httpServer *http.Server,
) {
	g.Go(func() error {
		go func() { //грейсфул остановка
			<-ctx.Done()
			ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), h.ShutdownTimeout)
			defer cancel()

			if err := httpServer.Shutdown(ctx); err != nil {
				logger(ctx).Error("http server shutdown error:", err)
			}
		}()

		logger(ctx).Info("http server started", slog.String("address", httpServer.Addr))

		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("httpServer.ListenAndServe: %w", err)
		}

		logger(ctx).Info("http server stopped", slog.String("address", httpServer.Addr))

		return nil
	})
}
