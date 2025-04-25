package restApp

import (
	"context"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	openapi "github.com/tummerk/golang/schedules/generatedOpenapi/go"
	"github.com/tummerk/golang/schedules/logger"
	"log/slog"
	"net/http"
	"time"
)

type Controller interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetUserSchedules(w http.ResponseWriter, r *http.Request)
	GetUserSchedule(w http.ResponseWriter, r *http.Request)
	NextTakings(w http.ResponseWriter, r *http.Request)
}

type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

type appRest struct {
	router     *mux.Router
	controller *Controller
	logger     Logger
}

func NewAppRest(c Controller, logger Logger) *appRest {
	router := openapi.NewRouter()
	router.HandleFunc("/schedule/create", c.Create)
	router.HandleFunc("/schedule", c.GetUserSchedule)
	router.HandleFunc("/schedules", c.GetUserSchedules)
	router.HandleFunc("/next_takings", c.NextTakings)

	app := appRest{
		router:     router,
		controller: &c,
		logger:     logger,
	}
	app.router.Use(app.MiddleWare)
	return &app
}

func (app *appRest) Run(addr string) {
	app.logger.Info("http server запустился", slog.String("addr", addr))
	e := http.ListenAndServe(addr, app.router)
	if e != nil {
		app.logger.Error(e.Error(), slog.String("addr", addr),
			slog.String("description", "http server лёг"))
		panic(e)
	}
}

func (app *appRest) MiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		traceID := uuid.NewString()

		ctx := context.WithValue(r.Context(), "traceID", traceID)
		ctx = context.WithValue(ctx, "startTime", startTime)
		ctx = context.WithValue(ctx, "ip", r.RemoteAddr)
		ctx = context.WithValue(ctx, "userAgent", r.UserAgent())
		r = r.WithContext(ctx)

		app.logger.Info("request", logger.RequestLog(r))

		lw := logger.LoggingResponseWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
			Size:           0,
		}
		next.ServeHTTP(&lw, r)
		app.logger.Info("response", logger.ResponseLog(ctx, lw))
	})
}
