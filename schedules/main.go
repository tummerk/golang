package main

import (
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {
	//роуты
	http.HandleFunc("/schedule/create", scheduleCreateHandler)
	http.HandleFunc("/schedule", scheduleHandler)
	http.HandleFunc("/schedules", userSchedulesHandler)
	http.HandleFunc("/next_takings", nextTakingsHandler)

	//запуск сервера
	http.ListenAndServe(":5252", nil)
}

// запуск миграций
func init() {

	e := RunMigrations()
	if e != nil {
		log.Fatalf("ошибка с миграциями %v", e)
	}
}
