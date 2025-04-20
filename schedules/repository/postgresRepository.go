package repository

import (
	"database/sql"
	"first_project/config"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"log"
	"time"
)

type PostgresRepository struct {
	DB *sql.DB
}

func (r *PostgresRepository) Connect() error {
	var e error
	r.DB, e = sql.Open("postgres", config.ConnStr)
	if e != nil {
		log.Fatal(e)
	}
	return nil
}

func (r *PostgresRepository) Close() {
	r.DB.Close()
}

func (r *PostgresRepository) GetUserSchedules(userID int) (Rows, error) {
	rows, e := r.DB.Query(`SELECT "id","medicament_name","receptions_per_day","date_start","date_end"
							   FROM schedules WHERE user_id = $1`, userID)
	if e != nil {
		log.Fatal(e)
	}
	fmt.Println(rows)
	return rows, e
}

func (r *PostgresRepository) GetUserSchedule(userID, scheduleID int) (Rows, error) {
	row, e := r.DB.Query(`SELECT "id","medicament_name","receptions_per_day","date_start","date_end"
							   FROM schedules WHERE user_id = $1 and id = $2`, userID, scheduleID)
	if e != nil {
		log.Fatal(e)
	}
	row.Next()
	return row, nil
}

func (r *PostgresRepository) NewUserSchedule(medicamentName string, userId, receptionsPerDay, duration int) (int, error) {
	tomorrow := time.Now().AddDate(0, 0, 1)
	dateStart := tomorrow.Format("2006-01-02")
	dateEnd := tomorrow.AddDate(0, 0, duration).Format("2006-01-02")

	result := r.DB.QueryRow(`INSERT INTO
								   schedules (
								     "medicament_name",
								     "user_id",
								     "receptions_per_day",
								     "date_start",
								     "date_end"
								   )
								 VALUES
	($1, $2, $3, $4, $5)
	RETURNING id`, medicamentName, userId, receptionsPerDay, dateStart, dateEnd)

	if result.Err() != nil {
		log.Fatal(result.Err())
	}

	var id int
	result.Scan(&id)
	return id, nil
}

func (r *PostgresRepository) RunMigrations() error {
	m, e := migrate.New("file://migrations", config.ConnStr)
	if e != nil {
		return e
	}
	defer m.Close()
	if e = m.Up(); e != nil && e != migrate.ErrNoChange {
		return e
	}
	return nil
}
