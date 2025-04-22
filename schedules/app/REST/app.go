package restApp

import (
	"github.com/gorilla/mux"
	openapi "github.com/tummerk/golang/schedules/generatedOpenapi/go"
	"net/http"
)

type Controller interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetUserSchedules(w http.ResponseWriter, r *http.Request)
	GetUserSchedule(w http.ResponseWriter, r *http.Request)
	NextTakings(w http.ResponseWriter, r *http.Request)
}

type appRest struct {
	router     *mux.Router
	controller *Controller
}

func NewAppRest(c Controller) *appRest {
	router := openapi.NewRouter()

	router.HandleFunc("/schedule/create", c.Create)
	router.HandleFunc("/schedule", c.GetUserSchedule)
	router.HandleFunc("/schedules", c.GetUserSchedules)
	router.HandleFunc("/next_takings", c.NextTakings)

	app := appRest{
		router:     router,
		controller: &c,
	}
	return &app
}

func (app *appRest) Run() {
	http.ListenAndServe(":5252", app.router)
}
