package app

import (
	"first_project/controller"
	"net/http"
)

type App struct {
	cont controller.ScheduleController
}

func NewApp(cont *controller.ScheduleController) *App {
	return &App{*cont}
}

func (a App) Run() {
	a.cont.UC.Repository.Connect()
	defer a.cont.UC.Repository.Close()

	http.HandleFunc("/schedule/create", a.cont.Create)
	http.HandleFunc("/schedule", a.cont.GetUserSchedule)
	http.HandleFunc("/schedules", a.cont.GetUserSchedules)
	http.HandleFunc("/next_takings", a.cont.NextTakings)

	//запуск сервера
	http.ListenAndServe(":5252", nil)
}
