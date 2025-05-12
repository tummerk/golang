package tests

import (
	"context"
	"database/sql"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/tummerk/golang/schedules/internal/application"
	serviceSchedule "github.com/tummerk/golang/schedules/internal/domain/service/schedule"
	"github.com/tummerk/golang/schedules/internal/infrastructure/persistance"
	grpcClient "github.com/tummerk/golang/schedules/internal/server/generated/grpc"
	"github.com/tummerk/golang/schedules/internal/server/rest"
	"github.com/tummerk/golang/schedules/pkg/application/modules"
	"github.com/tummerk/golang/schedules/pkg/middlewarex"
	"github.com/tummerk/golang/schedules/pkg/tests"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"net/http"
	"testing"
	"time"
)

type Suite struct {
	suite.Suite
	Container  *postgres.PostgresContainer
	DB         *sql.DB
	Service    serviceSchedule.ScheduleService
	httpServer *modules.HTTPServer
	grpcServer *modules.GrpcServer
	apiClient  tests.APIClient
	grpcClient grpcClient.ScheduleServiceClient
}

func TestIntegration(t *testing.T) {
	suite.Run(t, &Suite{})
}

func (s *Suite) SetupSuite() {
	var err error
	rq := s.Require()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Minute)

	s.Container, err = tests.NewPostgresContainer(ctx)
	rq.NoError(err, "Не удалось запустить контейнер PostgreSQL")

	connStr, err := s.Container.ConnectionString(ctx, "sslmode=disable")
	rq.NoError(err, "Не удалось получить строку подключения")

	s.DB, err = sql.Open("postgres", connStr)
	rq.NoError(err, "Не удалось подключиться к базе данных")

	err = s.DB.PingContext(ctx)
	rq.NoError(err, "Не удалось пингануть базу данных")

	s.apiClient = tests.NewAPIClient("http://localhost:5252", http.DefaultClient)

	conn, err := grpc.NewClient(":11111", grpc.WithTransportCredentials(insecure.NewCredentials()))

	s.grpcClient = grpcClient.NewScheduleServiceClient(conn)

	repo := persistance.NewRepo(s.DB)
	s.Service = serviceSchedule.NewScheduleService(&repo, 120)
	router := chi.NewRouter()
	router.Use(middlewarex.Logger, middlewarex.UserID([]byte("8a1f3d9c7b2e45f60a9e8d2b4c3fds76")))
	restServer.NewServer(&s.Service).RegisterRoutes(router)

	g := errgroup.Group{}
	httpServer := &http.Server{
		Addr:    ":5252",
		Handler: router,
	}
	s.httpServer = &modules.HTTPServer{ShutdownTimeout: 0}
	s.httpServer.Run(ctx, &g, httpServer)

	grpcListener, err := net.Listen("tcp", ":11111")
	s.grpcServer.Run(ctx, &g, application.NewGrpcServer(ctx, s.Service, []byte("8a1f3d9c7b2e45f60a9e8d2b4c3fds76")), grpcListener)
	time.Sleep(time.Second * 2)

}

func (s *Suite) TearDownSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.DB.Close()
	s.Container.Terminate(ctx)

}

func (s *Suite) SetupTest() {
	_, err := s.DB.Exec(`
        								TRUNCATE TABLE schedules 
        								RESTART IDENTITY CASCADE
									   `)
	s.Require().NoError(err)
}
