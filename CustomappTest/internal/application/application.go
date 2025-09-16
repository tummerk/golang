package application

import (
	"CustomappTest/internal/domain/service"
	httpServer "CustomappTest/internal/server/http"
	"CustomappTest/pkg/modules"
	"context"
	"github.com/go-chi/chi/v5"
	"golang.org/x/sync/errgroup"
	"net"
	"net/http"
)

func Run(port string, rtp float64) {
	ctx := context.Background()

	//создание генератора
	service, err := service.NewGenerator(rtp)
	if err != nil {
		panic(err)
	}
	g, _ := errgroup.WithContext(ctx)
	//запуск http сервера
	modules.HttpServerRun(ctx, g, newHTTPServer(ctx, *service, port))
	g.Wait()
}

func newHTTPServer(ctx context.Context, service service.Generator, Addr string) *http.Server {
	router := chi.NewRouter()

	httpServer.NewServer(&service).RegisterRoutes(router)
	return &http.Server{
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
		Addr:    Addr,
		Handler: router,
	}
}
