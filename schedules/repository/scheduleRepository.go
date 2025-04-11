package repository

type ScheduleRepository interface {
	Connect() error
	GetUserSchedules(userID int) (Rows, error)
	GetUserSchedule(userID, scheduleID int) (Rows, error)
	NewUserSchedule(medicamentName string, userId, receptionsPerDay, duration int) (int, error)
	Close()
}

type Rows interface {
	Scan(dest ...interface{}) error
	Next() bool
}
