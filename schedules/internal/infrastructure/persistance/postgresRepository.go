package persistance

import (
	"context"
	"database/sql"
	"github.com/tummerk/golang/schedules/internal/config"
	serviceSchedule "github.com/tummerk/golang/schedules/internal/domain/service/schedule"
	"github.com/tummerk/golang/schedules/pkg/contextx"
	"github.com/tummerk/golang/schedules/pkg/utils"
	"log/slog"
	"strconv"
	"time"
)

type Logger interface {
	Info(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

type ScheduleRepository struct {
	DB *sql.DB
}

func NewRepo(db *sql.DB) ScheduleRepository {
	return ScheduleRepository{
		db,
	}
}

func (r *ScheduleRepository) GetUserSchedules(ctx context.Context, userID int) (serviceSchedule.Rows, error) {
	traceID := slog.Any("traceID", contextx.TraceIDFromContext(ctx))
	rows, e := r.DB.Query(`SELECT "id","medicament_name","receptions_per_day","date_start","date_end"
							   FROM schedules WHERE user_id = $1`, userID)
	userIdEncode, _ := utils.Encrypt(strconv.Itoa(userID), config.Key)
	if e != nil {
		logger(ctx).Error("Error getting user schedules from postgres", traceID,
			slog.String("error", e.Error()),
			slog.String("userID", userIdEncode))
		return nil, e
	}
	logger(ctx).Info("Successfully got user schedules from postgres", traceID,
		slog.String("userID", userIdEncode))
	return rows, e
}

func (r *ScheduleRepository) GetUserSchedule(ctx context.Context, userID, scheduleID int) (serviceSchedule.Rows, error) {
	traceID := slog.Any("traceID", contextx.TraceIDFromContext(ctx))
	userIdEncode, _ := utils.Encrypt(strconv.Itoa(userID), config.Key)
	row, e := r.DB.Query(`SELECT "id","medicament_name","receptions_per_day","date_start","date_end"
							   FROM schedules WHERE user_id = $1 and id = $2`, userID, scheduleID)
	if e != nil {
		logger(ctx).Error("Error getting schedule from postgres", traceID, slog.String("error", e.Error()),
			slog.String("userID", userIdEncode),
			slog.Int("scheduleID", scheduleID))
		return nil, e
	}
	logger(ctx).Info("Successfully got schedule from postgres", traceID,
		slog.String("userID", userIdEncode),
		slog.Int("scheduleID", scheduleID))
	row.Next()
	return row, nil
}

func (r *ScheduleRepository) NewUserSchedule(ctx context.Context, medicamentName string, userId, receptionsPerDay, duration int) (int, error) {
	traceID := slog.Any("traceID", contextx.TraceIDFromContext(ctx))
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
		logger(ctx).Error("error with creating user schedule in postgres", traceID,
			slog.String("error", result.Err().Error()))
		return -1, result.Err()
	}
	var scheduleID int
	result.Scan(&scheduleID)
	logger(ctx).Info("Successfully created user schedule in postgres", traceID,
		slog.Int("scheduleId", scheduleID))
	return scheduleID, nil
}
