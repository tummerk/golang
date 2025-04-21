package main

import (
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/tummerk/golang/schedules/app/REST"
	grpcApp "github.com/tummerk/golang/schedules/app/gRPC"
	"github.com/tummerk/golang/schedules/controller"
	"github.com/tummerk/golang/schedules/repository"
	"github.com/tummerk/golang/schedules/useCase"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}

	repo := repository.PostgresRepository{}
	useCase := useCase.ScheduleUC{Repository: &repo}
	cont := controller.ScheduleController{UC: useCase}

	appRest := restApp.NewAppRest(&cont)
	appGrpc := grpcApp.NewApp("12345", useCase)

	repo.Connect()

	wg.Add(2)
	go appRest.Run()
	go appGrpc.Run()
	wg.Wait()
}
