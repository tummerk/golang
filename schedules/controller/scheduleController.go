package controller

import (
	"first_project/entities"
	"first_project/repository"
	"first_project/useCase"
	"first_project/utils"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type ScheduleController struct {
	UC useCase.ScheduleUC
}

func NewScheduleController(repo repository.ScheduleRepository) *ScheduleController {
	return &ScheduleController{*useCase.NewScheduleUC(repo)}
}

func (c ScheduleController) Create(w http.ResponseWriter, r *http.Request) {
	t, e := template.ParseFiles("templates/scheduleCreate.html")
	if e != nil {
		fmt.Println(e)
	}
	t.ExecuteTemplate(w, "scheduleCreate", nil)
}

func (c ScheduleController) GetUserSchedules(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "неизвестный метод", http.StatusMethodNotAllowed)
		return
	}

	funcmap := template.FuncMap{ //передача функции
		"MinuteToTime": utils.MinuteToTime,
		"TimeToDate":   utils.TimeToDate,
	}
	t, e := template.New("").Funcs(funcmap).ParseFiles("templates/userSchedules.html")
	if e != nil {
		log.Fatal(e)
	}

	query := r.URL.Query()
	userID, e := strconv.Atoi(query.Get("user_id"))
	if e != nil {
		http.Error(w, "укажите число!", http.StatusBadRequest)
		return
	}

	currentSchedules, e, pastSchedules := c.UC.GetUserSchedules(userID)

	if e != nil || len(currentSchedules)+len(pastSchedules) == 0 {
		http.Error(w, "Такому пользователю лекарства не назначались!", http.StatusBadRequest)
		return
	}

	data := struct {
		CurrentSchedules []entities.Schedule
		PastSchedules    []entities.Schedule
	}{currentSchedules,
		pastSchedules}
	t.ExecuteTemplate(w, "userSchedules", data)
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
		scheduleID, e := c.UC.Create(medicamentName, userID, receptionsPerDay, duration)

		//возвращаем ID нового расписания
		w.Write([]byte(strconv.Itoa(int(scheduleID))))
	case r.Method == http.MethodGet: //поиск schedule по ID и user_id

		funcmap := template.FuncMap{ //передача функции преобразования минут в время
			"MinuteToTime": utils.MinuteToTime,
			"TimeToDate":   utils.TimeToDate,
		}
		t, e := template.New("").Funcs(funcmap).ParseFiles("templates/userSchedule.html")

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

		schedule, e, isRelevant := c.UC.GetUserSchedule(userID, scheduleID)
		if e != nil {
			http.Error(w, "такого расписания не существует", http.StatusBadRequest)
			return
		}

		data := struct {
			Schedule   entities.Schedule
			IsRelevant bool
		}{
			Schedule:   schedule,
			IsRelevant: isRelevant,
		}

		t.ExecuteTemplate(w, "userSchedule", data)

	default:
		http.Error(w, "неизвестный метод", http.StatusMethodNotAllowed)
	}
}

func (c ScheduleController) NextTakings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "неизвестный метод", http.StatusMethodNotAllowed)
		return
	}

	t, e := template.ParseFiles("templates/nextTakings.html")

	query := r.URL.Query()
	userID, e := strconv.Atoi(query.Get("user_id"))
	if e != nil {
		http.Error(w, "укажите число!", http.StatusBadRequest)
		return
	}

	nextTakings, e := c.UC.NextTakings(userID)
	if e != nil {
		http.Error(w, "неверный user_id", http.StatusBadRequest)
		return
	}

	data := struct {
		NextTakings []string
	}{nextTakings}
	fmt.Println(data)
	t.ExecuteTemplate(w, "nextTakings", data)
}
