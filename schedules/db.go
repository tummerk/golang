package main

import (
	"database/sql"
	"first_project/config"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"time"
)

/*с базой данных было несколько проблем:
1: стоял вопрос необходимости таблицы для user, но т. к. авторизация/аутентификация не подразумевается я подумал что
 можно обойтись без неё ведь при создании новых расписаний для пользователей которых нет в бд приходилось бы создавать
 еще и их
2 (важно): хранить в базе данных таблицу уже с посчитанным временем приёмов и просто брать оттуда время приёма(допустим
при создании указываем что будет два приёма парацетамола и таблицу заносим:

название    время приемя   дата конца
парацетамол     8:00        29.05.25
парацетамол     22:00       29.05.25

либо же
хранить в базе данных только одну запись о расписании, в которой указано число приемов в день и просто считать каждый
раз время приёмов на день

название    приёмов в день   дата конца
парацетамол        2       	  29.05.25

я подумал что при большом количестве пользователей нагрузка на базу данных в 2 случае будет значительно меньше, и выбрал
его.


*/

var db *sql.DB

var connStr = "postgres://" + config.DB_USER + ":" + config.DB_PASS + "@" + config.DB_HOST + ":" + config.DB_PORT + "/" + config.DB_NAME + "?sslmode=disable"

// запуск миграций
func RunMigrations() error {
	m, e := migrate.New("file://migrations", connStr)
	if e != nil {
		return e
	}
	defer m.Close()
	if e = m.Up(); e != nil && e != migrate.ErrNoChange {
		return e
	}
	return nil
}

// подключение к бд
func connect() error {
	var e error
	db, e = sql.Open("postgres", connStr)
	if e != nil {
		return e
	}
	return nil
}

// новое расписание
func NewSchedule(medicamentName string, userId, receptionsPerDay, duration int) (int, error) {
	e := connect()
	if e != nil {
		return 0, e
	}
	defer db.Close()

	today := time.Now()
	tomorrow := today.AddDate(0, 0, 1)

	dateStart := tomorrow.Format("2006-01-02")
	dateEnd := tomorrow.AddDate(0, 0, duration).Format("2006-01-02")

	result := db.QueryRow(`INSERT INTO
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
	if e != nil {
		return 0, e
	}
	var id int
	result.Scan(&id)
	return id, nil
}

// все актуальные расписания приёма для человека
func UserShedules(userId int) ([]Schedule, error, []Schedule) {
	e := connect()
	if e != nil {
		return nil, e, nil
	}
	defer db.Close()

	currentSchedules := []Schedule{} //действительные
	pastSchedules := []Schedule{}    //которые уже прошли

	rows, e := db.Query(`SELECT "id","medicament_name","receptions_per_day","date_start","date_end"
							   FROM schedules WHERE user_id = $1`, userId)
	if e != nil {
		return nil, e, nil
	}
	var s Schedule
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

// определённое расписание для человека
func UserShedule(userId, scheduleId int) (Schedule, error, bool) {
	e := connect()
	if e != nil {
		return Schedule{}, e, false
	}
	defer db.Close()

	var s Schedule
	row := db.QueryRow(`SELECT "id","medicament_name","receptions_per_day","date_start","date_end"
							   FROM schedules WHERE user_id = $1 and id = $2`, userId, scheduleId)

	e = row.Scan(&s.ID, &s.MedicamentName, &s.ReceptionsPerDay, &s.DateStart, &s.DateEnd)
	if e != nil {
		return Schedule{}, e, false
	}
	isRelevant := false
	if time.Now().Before(s.DateEnd) {
		isRelevant = true
	}
	return s, nil, isRelevant
}

// ближайщие приёмы лекарств
func NextTakings(userId int) ([]string, error) {
	schedules, e, _ := UserShedules(userId)

	var nextTakings []string

	if e != nil {
		return nextTakings, e
	}

	minuteFromStartDay := MinuteFromStartDay(time.Now())
	for _, schedule := range schedules {
		for _, minute := range schedule.ScheduleOnDay() {
			switch {
			case minute < minuteFromStartDay:
				continue
			case minute-minuteFromStartDay < config.TIME_NEXT_TAKINGS:
				take := fmt.Sprintf("%s %s", MinuteToTime(minute), schedule.MedicamentName)
				nextTakings = append(nextTakings, take)
			default:
				break
			}
		}
	}
	return nextTakings, nil
}
