package grpcApp

import (
	scheduleGRPC "github.com/tummerk/golang/schedules/grpc"
	"github.com/tummerk/golang/schedules/useCase"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Debug(msg string, keysAndValues ...interface{})
}

type App struct {
	gRPCServer *grpc.Server
	port       string
	logger     Logger
}

func NewApp(port string, UC useCase.ScheduleUC, logger Logger) *App {
	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(scheduleGRPC.LoggingInterceptor(logger)))
	scheduleGRPC.Register(gRPCServer, UC, logger)
	return &App{gRPCServer, port, logger}
}

func (app *App) Run() error {
	l, err := net.Listen("tcp", ":"+app.port)
	if err != nil {
		app.logger.Error("Error listening on port "+app.port, slog.String("err", err.Error()))
		return err
	}
	app.logger.Info("gRPC сервер запустился", slog.String("port", app.port))
	err = app.gRPCServer.Serve(l)
	if err != nil {
		app.logger.Error("Error starting server", slog.String("err", err.Error()))
		return err
	}
	return nil
}

func (app *App) Stop() error {
	app.logger.Info("Stopping app server", slog.String("port", app.port))
	app.gRPCServer.GracefulStop()
	return nil
}
