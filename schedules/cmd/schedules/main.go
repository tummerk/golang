package main

import (
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/tummerk/golang/schedules/internal/application"
)

func main() {
	app := application.New()
	app.Run()
}
