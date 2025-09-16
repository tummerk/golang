package connectors

import (
	"context"
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"log/slog"
	"sync"
)

type Postgres struct {
	db   *sql.DB
	DSN  string
	init sync.Once
}

func (p *Postgres) Client(ctx context.Context) *sql.DB {
	var e error
	p.init.Do(func() {
		p.db, e = sql.Open("postgres", p.DSN)
		if e != nil {
			logger(ctx).Error("Error connecting to Postgres", slog.String("error", e.Error()))
		}
		e = p.db.PingContext(ctx)
		if e != nil {
			logger(ctx).Error("Error connecting to Postgres", slog.String("error", e.Error()))
		}
		logger(ctx).Info("Successfully connected to Postgres database")
	})
	return p.db
}

func (p *Postgres) Close(ctx context.Context) {
	if err := p.db.Close(); err != nil {
		logger(ctx).Error("postgresClient.Close", slog.String("error", err.Error()))
	}

	logger(ctx).Info(
		"postgres disconnected",
		slog.String("database", p.DSN),
	)
}

func (p *Postgres) RunMigrations(ctx context.Context) error {
	m, e := migrate.New("file://migrations", "postgres://postgres:pass123@db:5432/schedules?sslmode=disable")
	if e != nil {
		logger(ctx).Error("Error running migrations",
			slog.String("error", e.Error()))
		return e
	}
	defer m.Close()
	if e = m.Up(); e != nil {
		logger(ctx).Error("Error running migrations",
			slog.String("error", e.Error()))
		return e
	}
	logger(ctx).Info("Successfully migrated migrations")
	return nil
}
