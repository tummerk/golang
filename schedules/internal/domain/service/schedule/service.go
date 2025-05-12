package serviceSchedule

import (
	"context"
	"github.com/tummerk/golang/schedules/internal/domain/entity"
	"github.com/tummerk/golang/schedules/internal/domain/value"
	"github.com/tummerk/golang/schedules/pkg/contextx"
	"github.com/tummerk/golang/schedules/pkg/utils"
	"log/slog"
	"strconv"

	"time"
)

type ScheduleRepository interface {
	GetUserSchedules(ctx context.Context, userID int) (Rows, error)
	GetUserSchedule(ctx context.Context, userID, scheduleID int) (Rows, error)
	NewUserSchedule(ctx context.Context, medicamentName string, userId, receptionsPerDay, duration int) (int, error)
}

type Rows interface {
	Scan(dest ...interface{}) error
	Next() bool
}

type ScheduleService struct {
	Repository      ScheduleRepository
	TimeNextTakings int
}

func NewScheduleService(repository ScheduleRepository, timeNextTakings int) ScheduleService {
	return ScheduleService{
		Repository:      repository,
		TimeNextTakings: timeNextTakings,
	}
}

type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Debug(msg string, keysAndValues ...interface{})
}

func (service *ScheduleService) Create(ctx context.Context, medicamentName string, userId, receptionsPerDay, duration int) (int, error) {
	traceID := slog.Any("traceID", contextx.TraceIDFromContext(ctx))
	scheduleID, e := service.Repository.NewUserSchedule(ctx, medicamentName, userId, receptionsPerDay, duration)
	if e != nil {
		logger(ctx).Error("error while creating user schedule", traceID,
			slog.String("error", e.Error()))
		return 0, e
	}
	logger(ctx).Info("created user schedule", traceID, slog.Int("scheduleID", scheduleID))
	return scheduleID, nil
}

func (service *ScheduleService) GetUserSchedules(ctx context.Context) ([]entity.Schedule, error, []entity.Schedule) {
	traceID := slog.Any("traceID", contextx.TraceIDFromContext(ctx))
	userIDstr, e := contextx.UserIDFromContext(ctx)
	if e != nil {
		logger(ctx).Error("error getting userID", traceID,
			slog.String("error", e.Error()))
	}
	userID, _ := strconv.Atoi(userIDstr.String())
	maskTraceID := contextx.TraceIDFromContext(ctx)
	rows, e := service.Repository.GetUserSchedules(ctx, userID)
	if e != nil {
		logger(ctx).Error("error while creating user schedule", traceID,
			slog.String("error", e.Error()))
	}

	var currentSchedules []entity.Schedule //действительные
	var pastSchedules []entity.Schedule

	var s entity.Schedule
	for rows.Next() {
		e = rows.Scan(&s.ID, &s.MedicamentName, &s.ReceptionsPerDay, &s.DateStart, &s.DateEnd)
		if e != nil {
			return nil, e, nil
		}
		if time.Now().Before(s.DateEnd) {
			currentSchedules = append(currentSchedules, s)
		} else {
			pastSchedules = append(pastSchedules, s)
		}
	}
	logger(ctx).Info("get user schedules", traceID, "userID", maskTraceID)
	return currentSchedules, nil, pastSchedules
}

func (service *ScheduleService) GetUserSchedule(ctx context.Context, scheduleID int) (entity.Schedule, error, bool) {
	traceID := slog.Any("traceID", contextx.TraceIDFromContext(ctx))
	userIDstr, e := contextx.UserIDFromContext(ctx)
	if e != nil {
		logger(ctx).Error("error getting userID", traceID,
			slog.String("error", e.Error()))
	}
	userID, _ := strconv.Atoi(userIDstr.String())
	row, e := service.Repository.GetUserSchedule(ctx, userID, scheduleID)
	if e != nil {
		logger(ctx).Error("error while getting user schedule", traceID, slog.String("error", e.Error()))
		return entity.Schedule{}, e, false
	}

	var schedule entity.Schedule

	e = row.Scan(&schedule.ID, &schedule.MedicamentName, &schedule.ReceptionsPerDay, &schedule.DateStart, &schedule.DateEnd)
	if e != nil {
		logger(ctx).Error("error while getting user schedule", traceID, "error", e.Error())
		return entity.Schedule{}, e, false
	}

	isRelevant := time.Now().Before(schedule.DateEnd)
	logger(ctx).Info("get user schedule ", traceID, slog.Int("scheduleID", scheduleID))
	return schedule, nil, isRelevant
}

func (service *ScheduleService) NextTakings(ctx context.Context) ([]value.Taking, error) {
	traceID := slog.Any("traceID", contextx.TraceIDFromContext(ctx))
	maskUserId, e := contextx.MaskUserIDFromContext(ctx)
	if e != nil {
		return nil, e
	}
	userIDstr, e := contextx.UserIDFromContext(ctx)
	if e != nil {
		logger(ctx).Error("error getting userID", traceID,
			slog.String("error", e.Error()))
	}
	_, e = strconv.Atoi(userIDstr.String())
	if e != nil {
		logger(ctx).Error("error parsing userID from string", traceID)
	}
	schedules, e, _ := service.GetUserSchedules(ctx)

	var nextTakings []value.Taking

	if e != nil {
		logger(ctx).Error("error while getting user schedules", traceID)
		return nextTakings, e
	}

	minuteFromStartDay := utils.MinuteFromStartDay(time.Now())
	for _, schedule := range schedules {
		for _, minute := range schedule.ScheduleOnDay(ctx) {
			switch {
			case minute < minuteFromStartDay:
				continue
			case minute-minuteFromStartDay < service.TimeNextTakings:
				taking := value.Taking{
					Name: schedule.MedicamentName,
					Time: utils.MinuteToTime(minute),
				}
				nextTakings = append(nextTakings, taking)
			default:
				break
			}
		}
	}
	logger(ctx).Info("get user next takings", traceID, slog.String("user_id", maskUserId.String()))
	return nextTakings, nil
}
