package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func scheduleCreateHandler(w http.ResponseWriter, r *http.Request) {
	t, e := template.ParseFiles("templates/scheduleCreate.html")
	if e != nil {
		fmt.Println(e)
	}
	t.ExecuteTemplate(w, "scheduleCreate", nil)
}

func scheduleHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost:
		e := r.ParseForm()
		if e != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		medicamentName := r.FormValue("medicamentName")
		userID, e := strconv.Atoi(r.FormValue("userID"))
		receptionsPerDay, e := strconv.Atoi(r.FormValue("receptionsPerDay"))
		duration, e := strconv.Atoi(r.FormValue("duration"))
		if e != nil {
			http.Error(w, "вы указали не число в полях где это нужно", http.StatusBadRequest)
			return
		}
		scheduleID, e := NewSchedule(medicamentName, userID, receptionsPerDay, duration)
		w.Write([]byte(strconv.Itoa(int(scheduleID))))
	case r.Method == http.MethodGet:
		funcmap := template.FuncMap{ //передача функции преобразования минут в время
			"MinuteToTime": MinuteToTime,
			"TimeToDate":   TimeToDate,
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

		schedule, e, isRelevant := UserShedule(userID, scheduleID)
		if e != nil {
			http.Error(w, "такого расписания не существует", http.StatusBadRequest)
			return
		}

		data := struct {
			Schedule   Schedule
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

func userSchedulesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "неизвестный метод", http.StatusMethodNotAllowed)
		return
	}

	funcmap := template.FuncMap{ //передача функции преобразования минут в время
		"MinuteToTime": MinuteToTime,
		"TimeToDate":   TimeToDate,
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
	currentSchedules, e, pastSchedules := UserShedules(userID)
	fmt.Println(currentSchedules)
	if e != nil || len(currentSchedules)+len(pastSchedules) == 0 {
		http.Error(w, "Такому пользователю лекарства не назначались!", http.StatusBadRequest)
		return
	}

	data := struct {
		CurrentSchedules []Schedule
		PastSchedules    []Schedule
	}{currentSchedules,
		pastSchedules}
	t.ExecuteTemplate(w, "userSchedules", data)
}

func nextTakingsHandler(w http.ResponseWriter, r *http.Request) {
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

	nextTakings, e := NextTakings(userID)
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
