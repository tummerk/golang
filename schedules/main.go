package main

import (
	"first_project/app"
	"first_project/controller"
	"first_project/repository"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	repo := repository.PostgresRepository{}
	cont := controller.NewScheduleController(&repo)

	a := app.NewApp(cont)
	a.Run()
}
