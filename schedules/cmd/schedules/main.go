package main

import (
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/tummerk/golang/schedules/internal/application"
	"time"
	_ "time/tzdata"
)

func main() {

	fmt.Println(time.Now())
	app := application.New()
	app.Run()
}
