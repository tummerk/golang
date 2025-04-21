package restApp

import (
	"github.com/gorilla/mux"
	"github.com/tummerk/golang/schedules/controller"
	openapi "github.com/tummerk/golang/schedules/generatedOpenapi/go"
	"net/http"
)

type appRest struct {
	router     *mux.Router
	controller controller.ScheduleController
}

func NewAppRest(controller *controller.ScheduleController) *appRest {
	router := openapi.NewRouter()

	router.HandleFunc("/schedule/create", controller.Create)
	router.HandleFunc("/schedule", controller.GetUserSchedule)
	router.HandleFunc("/schedules", controller.GetUserSchedules)
	router.HandleFunc("/next_takings", controller.NextTakings)

	app := appRest{
		router:     router,
		controller: *controller,
	}
	return &app
}

func (app *appRest) Run() {
	http.ListenAndServe(":5252", app.router)
}
