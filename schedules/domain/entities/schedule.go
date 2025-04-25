package entities

import (
	"context"
	"github.com/tummerk/golang/schedules/utils"
	"log/slog"
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

type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
}

func (schedule Schedule) ScheduleOnDay(ctx context.Context, logger Logger) []int { //создание массива состоящего из времени приемов (в мин от начала дня)
	traceID := slog.Any("traceID", ctx.Value("traceID"))
	var takingsMedications []int
	minuteFromStartDay := 480
	minuteAvailable := 840
	if schedule.ReceptionsPerDay == 1 {
		return append(takingsMedications, 480)
	}

	step := minuteAvailable / (schedule.ReceptionsPerDay - 1)
	time := minuteFromStartDay

	for i := 0; i < schedule.ReceptionsPerDay; i++ {
		takingsMedications = append(takingsMedications, utils.RoundUp(time, 15)) //добавляем кратное 15 минутам время
		time += step
	}
	logger.Debug("creating schedule on day", traceID, slog.Int("ReceptionsPerDay", schedule.ReceptionsPerDay),
		slog.Any("takings", takingsMedications))
	return takingsMedications
}

func (schedule Schedule) ScheduleOnDayString(ctx context.Context, logger Logger) []string {
	takingsMedications := make([]string, 0)
	for _, v := range schedule.ScheduleOnDay(ctx, logger) {
		takingsMedications = append(takingsMedications, utils.MinuteToTime(v))
	}
	return takingsMedications
}
