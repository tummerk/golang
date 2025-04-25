package useCase

import (
	"context"
	"github.com/tummerk/golang/schedules/config"
	"github.com/tummerk/golang/schedules/domain/entities"
	"github.com/tummerk/golang/schedules/domain/repository"
	openapi "github.com/tummerk/golang/schedules/generatedOpenapi/go"
	"github.com/tummerk/golang/schedules/utils"
	"log/slog"

	"time"
)

type ScheduleUC struct {
	Repository repository.ScheduleRepository
	Logger     Logger
}

type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Debug(msg string, keysAndValues ...interface{})
}

func (uc *ScheduleUC) Create(ctx context.Context, medicamentName string, userId, receptionsPerDay, duration int) (int, error) {
	traceID := slog.Any("traceID", ctx.Value("traceID"))
	scheduleID, e := uc.Repository.NewUserSchedule(medicamentName, userId, receptionsPerDay, duration)
	if e != nil {
		uc.Logger.Error("error while creating user schedule", traceID,
			slog.String("error", e.Error()))
		return 0, e
	}
	uc.Logger.Info("created user schedule", traceID, slog.Int("scheduleID", scheduleID))
	return scheduleID, nil
}

func (uc *ScheduleUC) GetUserSchedules(ctx context.Context, userID int) ([]entities.Schedule, error, []entities.Schedule) {
	traceID := slog.Any("traceID", ctx.Value("traceID"))
	rows, e := uc.Repository.GetUserSchedules(userID)
	if e != nil {
		uc.Logger.Error("error while creating user schedule", traceID,
			slog.String("error", e.Error()))
	}

	var currentSchedules []entities.Schedule //действительные
	var pastSchedules []entities.Schedule

	var s entities.Schedule
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
	uc.Logger.Info("get user schedules", traceID, "userID", userID)
	return currentSchedules, nil, pastSchedules
}

func (uc *ScheduleUC) GetUserSchedule(ctx context.Context, userID, scheduleID int) (entities.Schedule, error, bool) {
	traceID := slog.Any("traceID", ctx.Value("traceID"))
	row, e := uc.Repository.GetUserSchedule(userID, scheduleID)
	if e != nil {
		uc.Logger.Error("error while getting user schedule", traceID, slog.String("error", e.Error()))
		return entities.Schedule{}, e, false
	}

	var schedule entities.Schedule

	e = row.Scan(&schedule.ID, &schedule.MedicamentName, &schedule.ReceptionsPerDay, &schedule.DateStart, &schedule.DateEnd)
	if e != nil {
		uc.Logger.Error("error while getting user schedule", traceID, "error", e.Error())
		return entities.Schedule{}, e, false
	}

	isRelevant := time.Now().Before(schedule.DateEnd)
	uc.Logger.Info("get user schedule ", traceID, slog.Int("scheduleID", scheduleID))
	return schedule, nil, isRelevant
}

func (uc *ScheduleUC) NextTakings(ctx context.Context, userID int) ([]openapi.Taking, error) {
	traceID := slog.Any("traceID", ctx.Value("traceID"))
	schedules, e, _ := uc.GetUserSchedules(ctx, userID)

	var nextTakings []openapi.Taking

	if e != nil {
		uc.Logger.Error("error while getting user schedules", traceID)
		return nextTakings, e
	}

	minuteFromStartDay := utils.MinuteFromStartDay(time.Now())
	for _, schedule := range schedules {
		for _, minute := range schedule.ScheduleOnDay(ctx, uc.Logger) {
			switch {
			case minute < minuteFromStartDay:
				continue
			case minute-minuteFromStartDay < config.TIME_NEXT_TAKINGS:
				taking := openapi.Taking{
					Name: schedule.MedicamentName,
					Time: utils.MinuteToTime(minute),
				}
				nextTakings = append(nextTakings, taking)
			default:
				break
			}
		}
	}
	uc.Logger.Info("get user next takings", traceID, slog.Int("user", userID))
	return nextTakings, nil
}
