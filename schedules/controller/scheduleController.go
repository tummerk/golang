package controller

import (
	"context"
	"encoding/json"
	"github.com/tummerk/golang/schedules/domain/entities"
	openapi "github.com/tummerk/golang/schedules/generatedOpenapi/go"
	"html/template"
	"net/http"
	"strconv"
)

type scheduleUC interface {
	Create(ctx context.Context, medicamentName string, userId, receptionsPerDay, duration int) (int, error)
	GetUserSchedules(ctx context.Context, userID int) ([]entities.Schedule, error, []entities.Schedule)
	GetUserSchedule(ctx context.Context, userID, scheduleID int) (entities.Schedule, error, bool)
	NextTakings(ctx context.Context, userID int) ([]openapi.Taking, error)
}

type Logger interface {
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}

type ScheduleController struct {
	UC     scheduleUC
	Logger Logger
}

func (c ScheduleController) Create(w http.ResponseWriter, r *http.Request) {

	t, e := template.ParseFiles("templates/scheduleCreate.html")
	if e != nil {
		c.Logger.Error(e.Error())
	}
	t.ExecuteTemplate(w, "scheduleCreate", nil)

}

func (c ScheduleController) GetUserSchedules(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "неизвестный метод", http.StatusMethodNotAllowed)

		return
	}

	query := r.URL.Query()
	userID, e := strconv.Atoi(query.Get("user_id"))
	if e != nil {
		http.Error(w, "укажите число!", http.StatusBadRequest)
		return
	}

	currentSchedules, e, pastSchedules := c.UC.GetUserSchedules(r.Context(), userID)

	if e != nil || len(currentSchedules)+len(pastSchedules) == 0 {
		http.Error(w, "Такому пользователю лекарства не назначались!", http.StatusBadRequest)
		return
	}
	currentSchedulesJson := []openapi.Schedule{}
	for _, schedule := range currentSchedules {
		takings := schedule.ScheduleOnDayString(r.Context(), c.Logger)
		currentSchedulesJson = append(currentSchedulesJson, openapi.Schedule{schedule.MedicamentName, takings})
	}

	pastSchedulesJson := []openapi.Schedule{}
	for _, schedule := range pastSchedules {
		takings := schedule.ScheduleOnDayString(r.Context(), c.Logger)
		currentSchedulesJson = append(currentSchedulesJson, openapi.Schedule{schedule.MedicamentName, takings})
	}

	response := map[string]interface{}{
		"current_schedules": currentSchedulesJson,
		"past_schedules":    pastSchedulesJson,
	}
	json.NewEncoder(w).Encode(response)
}

func (c ScheduleController) GetUserSchedule(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost: //занесение нового расписания в бд
		e := r.ParseForm()
		if e != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		//получаем данные
		medicamentName := r.FormValue("medicamentName")
		userID, e1 := strconv.Atoi(r.FormValue("userID"))
		receptionsPerDay, e2 := strconv.Atoi(r.FormValue("receptionsPerDay"))
		duration, e3 := strconv.Atoi(r.FormValue("duration"))
		if receptionsPerDay > 15 || receptionsPerDay < 1 {
			http.Error(w, "количество приёмов должно быть от 1 до 15", http.StatusBadRequest)
			return
		}
		if e1 != nil || e2 != nil || e3 != nil {
			http.Error(w, "вы указали не целое число в полях где это нужно", http.StatusBadRequest)
			return
		}
		scheduleID, e := c.UC.Create(r.Context(), medicamentName, userID, receptionsPerDay, duration)

		//возвращаем ID нового расписания
		w.Write([]byte(strconv.Itoa(scheduleID)))
	case r.Method == http.MethodGet: //поиск schedule по ID и user_id
		query := r.URL.Query()
		userID, e := strconv.Atoi(query.Get("user_id"))
		if e != nil {
			http.Error(w, "user_id указан неверно", http.StatusBadRequest)
			return
		}
		scheduleID, e := strconv.Atoi(query.Get("schedule_id"))
		if e != nil {
			http.Error(w, "schedule_id указан неверно", http.StatusBadRequest)
			return
		}

		schedule, e, isRelevant := c.UC.GetUserSchedule(r.Context(), userID, scheduleID)

		takings := schedule.ScheduleOnDayString(r.Context(), c.Logger)
		var scheduleJSON = openapi.Schedule{MedicamentName: schedule.MedicamentName, Takings: takings}

		response := map[string]interface{}{
			"scheduleJSON": scheduleJSON,
			"isRelevant":   isRelevant,
		}
		json.NewEncoder(w).Encode(response)
	default:
		http.Error(w, "неизвестный метод", http.StatusMethodNotAllowed)
	}
}

func (c ScheduleController) NextTakings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "неизвестный метод", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	userID, e := strconv.Atoi(query.Get("user_id"))
	if e != nil {
		http.Error(w, "укажите число!", http.StatusBadRequest)
		return
	}

	nextTakings, e := c.UC.NextTakings(r.Context(), userID)
	if e != nil {
		http.Error(w, "неверный user_id", http.StatusBadRequest)
		return
	}
	response := map[string]interface{}{
		"next_takings": nextTakings,
	}
	json.NewEncoder(w).Encode(response)
}
