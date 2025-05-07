package application

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/tummerk/golang/schedules/internal/config"
	serviceSchedule "github.com/tummerk/golang/schedules/internal/domain/service/schedule"
	"github.com/tummerk/golang/schedules/internal/infrastructure/persistance"
	Schedulegrpc "github.com/tummerk/golang/schedules/internal/server/grpc"
	"github.com/tummerk/golang/schedules/internal/server/rest"
	"github.com/tummerk/golang/schedules/pkg/application/connectors"
	"github.com/tummerk/golang/schedules/pkg/application/modules"
	"github.com/tummerk/golang/schedules/pkg/contextx"
	"github.com/tummerk/golang/schedules/pkg/interceptorx"
	"github.com/tummerk/golang/schedules/pkg/middlewarex"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os/signal"
	"syscall"
)

type App struct {
	cfg        config.Config
	slog       *connectors.Slog
	httpServer *modules.HTTPServer
	grpcServer *modules.GrpcServer
	service    serviceSchedule.ScheduleService
	postgres   *connectors.Postgres
	repo       persistance.ScheduleRepository
}

func New() *App {
	cfg, e := config.Load()
	if e != nil {
		log.Panic(e)
	}
	return &App{
		cfg:        cfg,
		grpcServer: &modules.GrpcServer{},
		slog: &connectors.Slog{
			Debug:    cfg.Logger.Debug,
			FileName: cfg.Logger.Filename,
		},
		httpServer: &modules.HTTPServer{
			ShutdownTimeout: cfg.Http.ShutdownTimeout,
		},
		postgres: &connectors.Postgres{
			DSN: cfg.Postgres.DSN,
		},
	}
}
func (app *App) ShutDown(ctx context.Context) {
	app.postgres.Close(ctx)
}

func (app *App) Run() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	ctx = contextx.WithLogger(ctx, app.slog.Logger())
	defer app.ShutDown(ctx)

	app.postgres.RunMigrations(ctx)

	defer stop()
	app.repo = persistance.NewRepo(app.postgres.Client(ctx))
	app.service = serviceSchedule.NewScheduleService(&app.repo, app.cfg.Service.TimeNextTakings)

	g, ctx := errgroup.WithContext(ctx)

	grpcListener, err := net.Listen("tcp", app.cfg.Grpc.Addr)
	if err != nil {
		app.slog.Logger().Error("Failed to listen for gRPC", slog.String("error", err.Error()),
			slog.String("Addr", app.cfg.Grpc.Addr))
	}
	app.httpServer.Run(ctx, g, app.newHTTPServer(ctx, app.service))
	app.grpcServer.Run(ctx, g, app.newGrpcServer(ctx, app.service), grpcListener)
	g.Wait()
}

func (app *App) newHTTPServer(ctx context.Context, service serviceSchedule.ScheduleService) *http.Server {
	router := chi.NewRouter()

	router.Use(middlewarex.TraceID, middlewarex.Logger)
	rest.NewServer(&service).RegisterRoutes(router)
	return &http.Server{
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
		Addr:    app.cfg.Http.ListenAddr,
		Handler: router,
	}
}

func (app *App) newGrpcServer(ctx context.Context, service serviceSchedule.ScheduleService) *grpc.Server {
	logger := contextx.LoggerFromContextOrDefault(ctx)
	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptorx.TraceIDInterceptor(), interceptorx.LoggingInterceptor(logger)))
	Schedulegrpc.Register(gRPCServer, &service)
	return gRPCServer
}
