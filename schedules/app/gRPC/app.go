package grpcApp

import (
	scheduleGRPC "github.com/tummerk/golang/schedules/grpc"
	"github.com/tummerk/golang/schedules/useCase"
	"google.golang.org/grpc"
	"log"
	"net"
)

type App struct {
	gRPCServer *grpc.Server
	port       string
}

func NewApp(port string, UC useCase.ScheduleUC) *App {
	gRPCServer := grpc.NewServer()
	scheduleGRPC.Register(gRPCServer, UC)
	return &App{gRPCServer, port}
}

func (app *App) Run() error {
	l, err := net.Listen("tcp", ":"+app.port)
	if err != nil {
		return err
	}
	log.Println("gRPC слушает на порту " + app.port)
	err = app.gRPCServer.Serve(l)
	if err != nil {
		return err
	}
	return nil
}

func (app *App) Stop() error {
	app.gRPCServer.GracefulStop()
	return nil
}
