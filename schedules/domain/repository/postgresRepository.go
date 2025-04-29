package repository

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/tummerk/golang/schedules/config"
	"github.com/tummerk/golang/schedules/utils"
	"log/slog"
	"strconv"
	"time"
)

type Logger interface {
	Info(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

type PostgresRepository struct {
	DB     *sql.DB
	Logger Logger
}

func (r *PostgresRepository) Connect() error {
	var e error
	r.DB, e = sql.Open("postgres", config.ConnStr)
	if e != nil {
		r.Logger.Error("Error connecting to Postgres database",
			slog.String("error", e.Error()))

	}
	r.Logger.Info("Successfully connected to Postgres database")
	return nil
}

func (r *PostgresRepository) Close() {
	r.DB.Close()
}

func (r *PostgresRepository) GetUserSchedules(userID int) (Rows, error) {
	rows, e := r.DB.Query(`SELECT "id","medicament_name","receptions_per_day","date_start","date_end"
							   FROM schedules WHERE user_id = $1`, userID)
	userIdEncode, _ := utils.Encrypt(strconv.Itoa(userID), config.Key)
	if e != nil {
		r.Logger.Error("Error getting user schedules",
			slog.String("error", e.Error()),
			slog.String("userID", userIdEncode))
	}
	r.Logger.Info("Successfully got user schedules",
		slog.String("userID", userIdEncode))
	return rows, e
}

func (r *PostgresRepository) GetUserSchedule(userID, scheduleID int) (Rows, error) {
	userIdEncode, _ := utils.Encrypt(strconv.Itoa(userID), config.Key)
	row, e := r.DB.Query(`SELECT "id","medicament_name","receptions_per_day","date_start","date_end"
							   FROM schedules WHERE user_id = $1 and id = $2`, userID, scheduleID)
	if e != nil {
		r.Logger.Error("Error getting schedule", slog.String("error", e.Error()),
			slog.String("userID", userIdEncode),
			slog.Int("scheduleID", scheduleID))
	}
	r.Logger.Info("Successfully got schedule",
		slog.String("userID", userIdEncode),
		slog.Int("scheduleID", scheduleID))
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
		r.Logger.Error("error with creating user schedule",
			slog.String("error", result.Err().Error()))
	}
	var scheduleID int
	result.Scan(&scheduleID)
	r.Logger.Info("Successfully created user schedule",
		slog.Int("scheduleId", scheduleID))
	return scheduleID, nil
}

func (r *PostgresRepository) RunMigrations() error {
	m, e := migrate.New("file://migrations", config.ConnStr)
	if e != nil {
		return e
	}
	defer m.Close()
	if e = m.Up(); e != nil && e != migrate.ErrNoChange {
		r.Logger.Error("Error running migrations",
			slog.String("error", e.Error()))
		return e
	}
	r.Logger.Info("Successfully migrated migrations")
	return nil
}
