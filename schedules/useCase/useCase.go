package useCase

import (
	"first_project/config"
	"first_project/entities"
	"first_project/repository"
	"first_project/utils"
	"fmt"
	"log"
	"time"
)

type ScheduleUC struct {
	Repository repository.ScheduleRepository
}

func NewScheduleUC(repository repository.ScheduleRepository) *ScheduleUC {
	return &ScheduleUC{Repository: repository}
}

func (uc *ScheduleUC) Create(medicamentName string, userId, receptionsPerDay, duration int) (int, error) {
	scheduleID, e := uc.Repository.NewUserSchedule(medicamentName, userId, receptionsPerDay, duration)
	if e != nil {
		return 0, e
	}

	return scheduleID, nil
}

func (uc *ScheduleUC) GetUserSchedules(userID int) ([]entities.Schedule, error, []entities.Schedule) {
	rows, e := uc.Repository.GetUserSchedules(userID)
	if e != nil {
		log.Fatal(e)
	}

	currentSchedules := []entities.Schedule{} //действительные
	pastSchedules := []entities.Schedule{}

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
	return currentSchedules, nil, pastSchedules
}

func (uc *ScheduleUC) GetUserSchedule(userID, scheduleID int) (entities.Schedule, error, bool) {

	row, e := uc.Repository.GetUserSchedule(userID, scheduleID)
	if e != nil {
		log.Fatal(e)
	}

	var schedule entities.Schedule

	e = row.Scan(&schedule.ID, &schedule.MedicamentName, &schedule.ReceptionsPerDay, &schedule.DateStart, &schedule.DateEnd)
	if e != nil {
		return entities.Schedule{}, e, false
	}

	isRelevant := time.Now().Before(schedule.DateEnd)

	return schedule, nil, isRelevant
}

func (uc *ScheduleUC) NextTakings(userID int) ([]string, error) {
	schedules, e, _ := uc.GetUserSchedules(userID)

	var nextTakings []string

	if e != nil {
		return nextTakings, e
	}

	minuteFromStartDay := utils.MinuteFromStartDay(time.Now())
	for _, schedule := range schedules {
		for _, minute := range schedule.ScheduleOnDay() {
			switch {
			case minute < minuteFromStartDay:
				continue
			case minute-minuteFromStartDay < config.TIME_NEXT_TAKINGS:
				take := fmt.Sprintf("%s %s", utils.MinuteToTime(minute), schedule.MedicamentName)
				nextTakings = append(nextTakings, take)
			default:
				break
			}
		}
	}
	return nextTakings, nil
}
