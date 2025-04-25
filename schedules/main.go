package main

import (
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	restApp "github.com/tummerk/golang/schedules/app/REST"
	grpcApp "github.com/tummerk/golang/schedules/app/gRPC"
	"github.com/tummerk/golang/schedules/controller"
	"github.com/tummerk/golang/schedules/domain/repository"
	logger "github.com/tummerk/golang/schedules/logger"
	"github.com/tummerk/golang/schedules/useCase"
	"log/slog"
	"sync"
)

func main() {
	loggerSchedule := logger.NewLogger("schedule.log", slog.LevelDebug)
	wg := sync.WaitGroup{}

	repo := repository.PostgresRepository{}
	useCase := useCase.ScheduleUC{Repository: &repo, Logger: loggerSchedule}
	//рест
	cont := controller.ScheduleController{UC: &useCase, Logger: loggerSchedule}
	appRest := restApp.NewAppRest(&cont, loggerSchedule)
	//gRPC
	appGrpc := grpcApp.NewApp("12345", useCase)

	repo.Connect()

	wg.Add(2)
	go appRest.Run(":5252")
	go appGrpc.Run()
	wg.Wait()
}
