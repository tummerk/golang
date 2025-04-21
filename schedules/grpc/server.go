package scheduleGRPC

import (
	"context"
	grpcGen "github.com/tummerk/golang/schedules/generatedProtobuff/gen"
	"github.com/tummerk/golang/schedules/useCase"
	"google.golang.org/grpc"
	"log"
)

type serverAPI struct {
	grpcGen.UnimplementedScheduleServiceServer
	UC useCase.ScheduleUC
}

func Register(gRPC *grpc.Server, UC useCase.ScheduleUC) {
	grpcGen.RegisterScheduleServiceServer(gRPC, &serverAPI{UC: UC})
}

func (s *serverAPI) GetSchedule(ctx context.Context, req *grpcGen.GetScheduleRequest) (*grpcGen.Schedule, error) {
	schedule, e, isRelevant := s.UC.GetUserSchedule(int(req.GetUserID()), int(req.GetScheduleID()))
	if e != nil {
		log.Print(e)
	}

	return &grpcGen.Schedule{
		MedicamentName: schedule.MedicamentName,
		IsActual:       isRelevant,
		Takings:        schedule.ScheduleOnDayString(),
	}, nil
}

func (s *serverAPI) GetSchedules(ctx context.Context, req *grpcGen.UserID) (*grpcGen.Schedules, error) {
	currentSchedules, e, _ := s.UC.GetUserSchedules(int(req.GetUserID()))
	if e != nil {
		log.Print(e)
	}
	var schedules []*grpcGen.Schedule
	for _, schedule := range currentSchedules {
		schedules = append(schedules, &grpcGen.Schedule{
			MedicamentName: schedule.MedicamentName,
			IsActual:       true,
			Takings:        schedule.ScheduleOnDayString(),
		})
	}

	return &grpcGen.Schedules{
		CurrentSchedules: schedules,
	}, nil
}

func (s *serverAPI) CreateSchedule(ctx context.Context, req *grpcGen.CreateScheduleRequest) (*grpcGen.ScheduleID, error) {
	userID := int(req.UserId)
	medicamentName := req.MedicamentName
	receptionsPerDay := int(req.ReceptionsPerDay)
	duration := int(req.Duration)

	if userID == 0 || medicamentName == "" || receptionsPerDay == 0 || duration == 0 {
		return &grpcGen.ScheduleID{ScheduleID: -1}, nil
	}

	scheduleID, e := s.UC.Create(medicamentName, userID, receptionsPerDay, duration)

	if e != nil {
		log.Print(e)
	}
	return &grpcGen.ScheduleID{ScheduleID: int64(scheduleID)}, nil
}

func (s *serverAPI) NextTakings(ctx context.Context, req *grpcGen.UserID) (*grpcGen.Takings, error) {
	takings, e := s.UC.NextTakings(int(req.UserID))
	if e != nil {
		log.Print(e)
	}
	var takingsGRPC grpcGen.Takings
	for _, v := range takings {
		takingsGRPC.Takings = append(takingsGRPC.Takings, &grpcGen.Taking{
			Name: v.Name,
			Time: v.Time,
		})
	}
	return &takingsGRPC, nil
}
