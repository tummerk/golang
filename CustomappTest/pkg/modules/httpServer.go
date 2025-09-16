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

func HttpServerRun(
	ctx context.Context,
	g *errgroup.Group,
	httpServer *http.Server,
) {
	g.Go(func() error {
		go func() { //грейсфул остановка
			<-ctx.Done()
			var cancel context.CancelFunc
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := httpServer.Shutdown(ctx); err != nil {
				slog.Error("http server shutdown error:", slog.String("error", err.Error()))
			}
		}()

		slog.Info("http server started", slog.String("address", httpServer.Addr))

		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("httpServer.ListenAndServe: %w", err)
		}

		slog.Info("http server stopped", slog.String("address", httpServer.Addr))

		return nil
	})
}
