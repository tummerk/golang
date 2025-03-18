package main

import (
	"time"
)

type Schedule struct {
	ID               int
	MedicamentName   string
	ReceptionsPerDay int
	DateStart        time.Time
	DateEnd          time.Time
}

func (schedule Schedule) ScheduleOnDay() []int {
	var takingMedications []int
	minuteFromStartDay := 480
	minuteAvaliable := 840
	if schedule.ReceptionsPerDay == 1 {
		return append(takingMedications, 480)
	}

	step := minuteAvaliable / (schedule.ReceptionsPerDay - 1)
	time := minuteFromStartDay

	for i := 0; i < schedule.ReceptionsPerDay; i++ {
		takingMedications = append(takingMedications, roundUp(time, 15)) //добавляем кратное 15 минутам время
		time += step
	}

	return takingMedications
}
