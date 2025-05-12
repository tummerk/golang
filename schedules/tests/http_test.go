package tests

import (
	"bou.ke/monkey"
	"context"
	"github.com/tummerk/golang/schedules/internal/domain/entity"
	"github.com/tummerk/golang/schedules/pkg/rest"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type ScheduleCreateResponse struct {
	ID int `json:"id"`
}

type ScheduleGetResponse struct {
	IsRelevant bool          `json:"isRelevant"`
	Schedule   rest.Schedule `json:"schedule"`
}

type SchedulesGetResponse struct {
	CurrentSchedules []rest.Schedule `json:"current_schedules"`
	PastSchedules    []rest.Schedule `json:"past_schedules"`
}

type NextTakingsResponse struct {
	Takings []rest.Taking `json:"next_takings"`
}

func (s *Suite) TestScheduleCreateHTTP() {
	form := url.Values{}
	form.Add("medicamentName", "test")
	form.Add("userID", "1")
	form.Add("receptionsPerDay", "15")
	form.Add("duration", "20")
	testCase := []struct {
		name             string
		bootstrap        func()
		expectedResponse ScheduleCreateResponse
		expectedStatus   int
		form             url.Values
	}{
		{
			name:             "success",
			bootstrap:        func() {},
			expectedResponse: ScheduleCreateResponse{ID: 1},
			expectedStatus:   http.StatusOK,
			form:             form,
		},
	}
	rq := s.Require()
	for _, tc := range testCase {
		s.Run(tc.name, func() {
			if tc.bootstrap != nil {
				tc.bootstrap()
			}
			var responseBody ScheduleCreateResponse
			res, err := s.apiClient.Post(context.Background(), "/schedule", nil, form, &responseBody, nil)
			rq.NoError(err)
			defer res.Body.Close()

			rq.Equal(tc.expectedStatus, res.StatusCode)
			rq.Equal(tc.expectedResponse, responseBody)

			var schedule entity.Schedule

			Row := s.DB.QueryRow(`SELECT medicament_name, user_id,receptions_per_day FROM schedules WHERE id = $1`, responseBody.ID)
			Row.Scan(&schedule.MedicamentName, &schedule.UserID, &schedule.ReceptionsPerDay)

			rq.Equal(form.Get("medicamentName"), schedule.MedicamentName)
			rq.Equal(form.Get("userID"), strconv.Itoa(schedule.UserID))
			rq.Equal(form.Get("receptionsPerDay"), strconv.Itoa(schedule.ReceptionsPerDay))

			s.Require().NoError(err)

		})
	}
	return
}

func (s *Suite) TestScheduleGetHTTP() {
	rq := s.Require()
	testCase := []struct {
		name             string
		bootstrap        func()
		expectedResponse ScheduleGetResponse
		expectedStatus   int
	}{
		{
			name: "success",
			bootstrap: func() {
				_, err := s.DB.Exec(`INSERT INTO schedules (medicament_name,user_id,receptions_per_day,date_start,date_end)
							values  ($1, $2, $3, $4, $5)`, "test", 1, 15, time.Now(), time.Now().Add(time.Hour*72))
				rq.NoError(err)
			},
			expectedResponse: ScheduleGetResponse{
				IsRelevant: true,
				Schedule: rest.Schedule{
					MedicamentName: "test",
					Takings: []string{"08:00", "09:00", "10:00", "11:00", "12:00", "13:00", "14:00", "15:00",
						"16:00", "17:00", "18:00", "19:00", "20:00", "21:00", "22:00"},
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCase {
		s.Run(tc.name, func() {
			if tc.bootstrap != nil {
				tc.bootstrap()
			}
			var responseBody ScheduleGetResponse

			res, err := s.apiClient.Get(context.Background(), "/schedule?user_id=1&schedule_id=1", nil, &responseBody, nil)
			rq.NoError(err)
			rq.Equal(tc.expectedStatus, res.StatusCode)
			rq.Equal(tc.expectedResponse, responseBody)
		})
	}
}

func (s *Suite) TestSchedulesGetHTTP() {
	rq := s.Require()
	testCase := []struct {
		name             string
		bootstrap        func()
		expectedResponse SchedulesGetResponse
		expectedStatus   int
		URL              string
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
			expectedResponse: SchedulesGetResponse{
				CurrentSchedules: []rest.Schedule{{
					MedicamentName: "test1",
					Takings: []string{"08:00", "09:00", "10:00", "11:00", "12:00", "13:00", "14:00", "15:00", "16:00",
						"17:00", "18:00", "19:00", "20:00", "21:00", "22:00"}},
				},
				PastSchedules: []rest.Schedule{{
					MedicamentName: "test2",
					Takings:        []string{"08:00"},
				}},
			},
			expectedStatus: http.StatusOK,
			URL:            "/schedules?user_id=1",
		},
	}
	for _, tc := range testCase {
		s.Run(tc.name, func() {
			if tc.bootstrap != nil {
				tc.bootstrap()
			}
			var responseBody SchedulesGetResponse
			res, err := s.apiClient.Get(context.Background(), tc.URL, nil, &responseBody, nil)
			rq.NoError(err)
			rq.Equal(tc.expectedStatus, res.StatusCode)
			rq.Equal(tc.expectedResponse, responseBody)
		})
	}
}

func (s *Suite) TestNextTakingsHTTP() {
	patch := monkey.Patch(time.Now, func() time.Time { return time.Date(2025, 5, 12, 7, 0, 0, 0, time.UTC) })
	defer patch.Unpatch()
	rq := s.Require()
	testCases := []struct {
		name             string
		bootstrap        func()
		expectedResponse NextTakingsResponse
		expectedStatus   int
		URL              string
	}{
		{
			name: "success",
			bootstrap: func() {
				_, err := s.DB.Exec(`INSERT INTO schedules (medicament_name,user_id,receptions_per_day,date_start,date_end)
							values  ($1, $2, $3, $4, $5)`, "test", 1, 2, time.Now().Add(-time.Hour*72),
					time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC))
				rq.NoError(err)
			},
			expectedResponse: NextTakingsResponse{
				Takings: []rest.Taking{{Name: "test", Time: "08:00"}},
			},
			expectedStatus: http.StatusOK,
			URL:            "/next_takings?user_id=1",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if tc.bootstrap != nil {
				tc.bootstrap()
			}
			var responseBody NextTakingsResponse
			res, err := s.apiClient.Get(context.Background(), tc.URL, nil, &responseBody, nil)
			rq.NoError(err)
			rq.Equal(tc.expectedStatus, res.StatusCode)
			rq.Equal(tc.expectedResponse, responseBody)
		})
	}
}
