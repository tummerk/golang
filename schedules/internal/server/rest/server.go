package restServer

import (
	"context"
	"encoding/json"
	"github.com/tummerk/golang/schedules/internal/domain/entity"
	"github.com/tummerk/golang/schedules/internal/domain/value"
	"github.com/tummerk/golang/schedules/pkg/contextx"
	"github.com/tummerk/golang/schedules/pkg/rest"
	"net/http"
	"strconv"
)

type ScheduleService interface {
	Create(ctx context.Context, medicamentName string, userId, receptionsPerDay, duration int) (int, error)
	GetUserSchedules(ctx context.Context) ([]entity.Schedule, error, []entity.Schedule)
	GetUserSchedule(ctx context.Context, scheduleID int) (entity.Schedule, error, bool)
	NextTakings(ctx context.Context) ([]value.Taking, error)
}

type Server struct {
	Service ScheduleService
}

func NewServer(s ScheduleService) *Server {
	return &Server{
		Service: s,
	}
}

func (s Server) GetUserSchedules(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "неизвестный метод", http.StatusMethodNotAllowed)

		return
	}

	userID, e := contextx.UserIDFromContext(r.Context())
	_, e = strconv.Atoi(userID.String())
	if e != nil {
		http.Error(w, "укажите число!", http.StatusBadRequest)
		return
	}

	currentSchedules, e, pastSchedules := s.Service.GetUserSchedules(r.Context())

	if e != nil || len(currentSchedules)+len(pastSchedules) == 0 {
		http.Error(w, "Такому пользователю лекарства не назначались!", http.StatusBadRequest)
		return
	}
	currentSchedulesJson := []rest.Schedule{}
	for _, schedule := range currentSchedules {
		takings := schedule.ScheduleOnDayString(r.Context())
		currentSchedulesJson = append(currentSchedulesJson, rest.Schedule{MedicamentName: schedule.MedicamentName, Takings: takings})
	}

	pastSchedulesJson := []rest.Schedule{}
	for _, schedule := range pastSchedules {
		takings := schedule.ScheduleOnDayString(r.Context())
		pastSchedulesJson = append(pastSchedulesJson, rest.Schedule{MedicamentName: schedule.MedicamentName, Takings: takings})
	}

	response := map[string]interface{}{
		"current_schedules": currentSchedulesJson,
		"past_schedules":    pastSchedulesJson,
	}
	json.NewEncoder(w).Encode(response)
}

func (s Server) CreateUserSchedule(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost { //занесение нового расписания в бд
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
		scheduleID, e := s.Service.Create(r.Context(), medicamentName, userID, receptionsPerDay, duration)

		//возвращаем ID нового расписания
		response := struct {
			ID int `json:"id"`
		}{ID: scheduleID}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "неизвестный метод", http.StatusMethodNotAllowed)
	}
}
func (s Server) GetUserSchedule(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet { //поиск schedule по ID и user_id
		query := r.URL.Query()
		_, e := contextx.UserIDFromContext(r.Context())
		if e != nil {
			http.Error(w, "user_id указан неверно", http.StatusBadRequest)
			return
		}
		scheduleID, e := strconv.Atoi(query.Get("schedule_id"))
		if e != nil {
			http.Error(w, "schedule_id указан неверно", http.StatusBadRequest)
			return
		}
		schedule, e, isRelevant := s.Service.GetUserSchedule(r.Context(), scheduleID)

		takings := schedule.ScheduleOnDayString(r.Context())
		var scheduleJSON = rest.Schedule{MedicamentName: schedule.MedicamentName, Takings: takings}

		response := map[string]interface{}{
			"schedule":   scheduleJSON,
			"isRelevant": isRelevant,
		}
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "неизвестный метод", http.StatusMethodNotAllowed)
	}
}

func (s Server) NextTakings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "неизвестный метод", http.StatusMethodNotAllowed)
		return
	}

	userID, e := contextx.UserIDFromContext(r.Context())
	_, e = strconv.Atoi(userID.String())
	if e != nil {
		http.Error(w, "укажите число!", http.StatusBadRequest)
		return
	}

	nextTakings, e := s.Service.NextTakings(r.Context())
	if e != nil {
		http.Error(w, "неверный user_id", http.StatusBadRequest)
		return
	}
	nextTakingsJson := []rest.Taking{}
	for _, t := range nextTakings {
		nextTakingsJson = append(nextTakingsJson, rest.Taking{
			Name: t.Name,
			Time: t.Time,
		})
	}
	response := map[string]interface{}{
		"next_takings": nextTakingsJson,
	}
	json.NewEncoder(w).Encode(response)
}
