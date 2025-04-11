package entities

import (
	"first_project/utils"
	"time"
)

/*
Почему в Schedule нет поля userID?
с текущим функционалом все команды подразумевают изначальное наличие userID в самом запросе
/schedules?user_id=2 ну и мы ищем schedule в базе данных только по user_id, конечно можно было бы добавить но это нигде
бы не использовалось
*/

type Schedule struct {
	ID               int
	MedicamentName   string
	ReceptionsPerDay int
	DateStart        time.Time
	DateEnd          time.Time
}

func (schedule Schedule) ScheduleOnDay() []int { //создание массива состоящего из времени приемов (в мин от начала дня)
	var takingMedications []int
	minuteFromStartDay := 480
	minuteAvaliable := 840
	if schedule.ReceptionsPerDay == 1 {
		return append(takingMedications, 480)
	}

	step := minuteAvaliable / (schedule.ReceptionsPerDay - 1)
	time := minuteFromStartDay

	for i := 0; i < schedule.ReceptionsPerDay; i++ {
		takingMedications = append(takingMedications, utils.RoundUp(time, 15)) //добавляем кратное 15 минутам время
		time += step
	}

	return takingMedications
}
