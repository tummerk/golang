package tests

import (
	"bou.ke/monkey"
	"context"
	"github.com/tummerk/golang/schedules/internal/domain/entity"
	grpcGen "github.com/tummerk/golang/schedules/internal/server/generated/grpc"
	"time"
)

func (s *Suite) TestScheduleCreateGRPC() {
	testCase := []struct {
		name             string
		bootstrap        func()
		expectedResponse grpcGen.ScheduleID
	}{
		{
			name:      "success",
			bootstrap: nil,
			expectedResponse: grpcGen.ScheduleID{
				ScheduleID: 1,
			},
		},
	}
	rq := s.Require()
	for _, tc := range testCase {
		s.Run(tc.name, func() {
			if tc.bootstrap != nil {
				tc.bootstrap()
			}
			req := grpcGen.CreateScheduleRequest{
				MedicamentName:   "test",
				UserId:           1,
				ReceptionsPerDay: 1,
				Duration:         1,
			}
			res, err := s.grpcClient.CreateSchedule(context.Background(), &req)
			rq.NoError(err)
			rq.Equal(tc.expectedResponse.ScheduleID, res.ScheduleID)
			var schedule entity.Schedule

			Row := s.DB.QueryRow(`SELECT medicament_name, user_id,receptions_per_day FROM schedules WHERE id = 1`)
			err = Row.Scan(&schedule.MedicamentName, &schedule.UserID, &schedule.ReceptionsPerDay)
			rq.NoError(err)
			rq.Equal(req.MedicamentName, schedule.MedicamentName)
			rq.Equal(int(req.UserId), schedule.UserID)
			rq.Equal(int(req.ReceptionsPerDay), schedule.ReceptionsPerDay)
		})
	}
}

func (s *Suite) TestScheduleGetGRPC() {
	rq := s.Require()
	testCase := []struct {
		name             string
		bootstrap        func()
		expectedResponse grpcGen.Schedule
	}{
		{
			name: "success",
			bootstrap: func() {
				_, err := s.DB.Exec(`INSERT INTO schedules (medicament_name,user_id,receptions_per_day,date_start,date_end)
							values  ($1, $2, $3, $4, $5)`, "test", 1, 15, time.Now(), time.Now().Add(time.Hour*72))
				rq.NoError(err)
			},
			expectedResponse: grpcGen.Schedule{
				MedicamentName: "test",
				Takings: []string{"08:00", "09:00", "10:00", "11:00", "12:00", "13:00", "14:00", "15:00",
					"16:00", "17:00", "18:00", "19:00", "20:00", "21:00", "22:00"},
			},
		},
	}
	for _, tc := range testCase {
		s.Run(tc.name, func() {
			if tc.bootstrap != nil {
				tc.bootstrap()
			}
			req := grpcGen.GetScheduleRequest{
				ScheduleID: 1,
				UserID:     1,
			}
			res, err := s.grpcClient.GetSchedule(context.Background(), &req)
			rq.NoError(err)
			rq.Equal(tc.expectedResponse.MedicamentName, res.MedicamentName)
			rq.Equal(tc.expectedResponse.Takings, res.Takings)
		})
	}
}

func (s *Suite) TestSchedulesGetGRPC() {
	rq := s.Require()
	testCase := []struct {
		name             string
		bootstrap        func()
		expectedResponse grpcGen.Schedules
	}{
		{
			name: "success",
			bootstrap: func() {
				_, err := s.DB.Exec(`INSERT INTO schedules (medicament_name,user_id,receptions_per_day,date_start,date_end)
							values  ($1, $2, $3, $4, $5)`, "test1", 1, 15, time.Now(), time.Now().Add(time.Hour*72))
				rq.NoError(err)
				_, err = s.DB.Exec(`INSERT INTO schedules (medicament_name,user_id,receptions_per_day,date_start,date_end)
							values  ($1, $2, $3, $4, $5)`, "test2", 1, 1, time.Now().Add(-time.Hour*72), time.Now().Add(-time.Hour*48))
				rq.NoError(err)
			},
			expectedResponse: grpcGen.Schedules{
				CurrentSchedules: []*grpcGen.Schedule{{
					MedicamentName: "test1",
					IsActual:       true,
					Takings: []string{"08:00", "09:00", "10:00", "11:00", "12:00", "13:00", "14:00", "15:00",
						"16:00", "17:00", "18:00", "19:00", "20:00", "21:00", "22:00"}},
				},
			},
		},
	}
	for _, tc := range testCase {
		s.Run(tc.name, func() {
			if tc.bootstrap != nil {
				tc.bootstrap()
			}
			req := grpcGen.UserID{UserID: 1}
			res, err := s.grpcClient.GetSchedules(context.Background(), &req)
			rq.NoError(err)
			rq.Equal(tc.expectedResponse.CurrentSchedules, res.CurrentSchedules)
		})
	}
}

func (s *Suite) TestNextTakingsGRPC() {
	patch := monkey.Patch(time.Now, func() time.Time { return time.Date(2025, 5, 12, 7, 0, 0, 0, time.UTC) })
	defer patch.Unpatch()
	rq := s.Require()
	testCase := []struct {
		name             string
		bootstrap        func()
		expectedResponse grpcGen.Takings
	}{
		{
			name: "success",
			bootstrap: func() {
				_, err := s.DB.Exec(`INSERT INTO schedules (medicament_name,user_id,receptions_per_day,date_start,date_end)
							values  ($1, $2, $3, $4, $5)`, "test", 1, 2, time.Now().Add(-time.Hour*72),
					time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC))
				rq.NoError(err)
			},
			expectedResponse: grpcGen.Takings{Takings: []*grpcGen.Taking{{Name: "test", Time: "08:00"}}},
		},
	}
	for _, tc := range testCase {
		s.Run(tc.name, func() {
			if tc.bootstrap != nil {
				tc.bootstrap()
			}
			req := grpcGen.UserID{UserID: 1}
			res, err := s.grpcClient.NextTakings(context.Background(), &req)
			rq.NoError(err)
			rq.Equal(tc.expectedResponse.Takings, res.Takings)
		})
	}
}
