package grpcServer

import (
	"context"
	"github.com/tummerk/golang/schedules/internal/domain/entity"
	_ "github.com/tummerk/golang/schedules/internal/domain/entity"
	"github.com/tummerk/golang/schedules/internal/domain/value"
	grpcGen "github.com/tummerk/golang/schedules/internal/server/generated/grpc"
	"google.golang.org/grpc"
	"log"
)

type ScheduleService interface {
	Create(ctx context.Context, medicamentName string, userId, receptionsPerDay, duration int) (int, error)
	GetUserSchedules(ctx context.Context) ([]entity.Schedule, error, []entity.Schedule)
	GetUserSchedule(ctx context.Context, scheduleID int) (entity.Schedule, error, bool)
	NextTakings(ctx context.Context) ([]value.Taking, error)
}

type serverAPI struct {
	grpcGen.UnimplementedScheduleServiceServer
	Service ScheduleService
}

func Register(gRPC *grpc.Server, Service ScheduleService) {
	grpcGen.RegisterScheduleServiceServer(gRPC, &serverAPI{Service: Service})
}

func (s *serverAPI) GetSchedule(ctx context.Context, req *grpcGen.GetScheduleRequest) (*grpcGen.Schedule, error) {
	schedule, e, isRelevant := s.Service.GetUserSchedule(ctx, int(req.GetScheduleID()))
	if e != nil {
		log.Print(e)
	}

	return &grpcGen.Schedule{
		MedicamentName: schedule.MedicamentName,
		IsActual:       isRelevant,
		Takings:        schedule.ScheduleOnDayString(ctx),
	}, nil
}

func (s *serverAPI) GetSchedules(ctx context.Context, req *grpcGen.UserID) (*grpcGen.Schedules, error) {
	currentSchedules, e, _ := s.Service.GetUserSchedules(ctx)
	if e != nil {
		log.Print(e)
	}
	var schedules []*grpcGen.Schedule
	for _, schedule := range currentSchedules {
		schedules = append(schedules, &grpcGen.Schedule{
			MedicamentName: schedule.MedicamentName,
			IsActual:       true,
			Takings:        schedule.ScheduleOnDayString(ctx),
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

	scheduleID, e := s.Service.Create(ctx, medicamentName, userID, receptionsPerDay, duration)

	if e != nil {
		log.Print(e)
	}
	return &grpcGen.ScheduleID{ScheduleID: int64(scheduleID)}, nil
}

func (s *serverAPI) NextTakings(ctx context.Context, req *grpcGen.UserID) (*grpcGen.Takings, error) {
	takings, e := s.Service.NextTakings(ctx)
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
