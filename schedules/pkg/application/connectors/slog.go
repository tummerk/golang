package connectors

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"log/slog"
	"sync"
)

type Slog struct {
	logger   *slog.Logger
	init     sync.Once
	Debug    bool
	FileName string
}

func (s *Slog) Logger() *slog.Logger {
	s.init.Do(func() {
		logRotator := &lumberjack.Logger{
			Filename:   s.FileName,
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     30,
		}

		var handler *slog.JSONHandler
		if s.Debug {
			handler = slog.NewJSONHandler(logRotator, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			})
			s.logger = slog.New(handler)
			return
		}
		handler = slog.NewJSONHandler(logRotator, nil)
		s.logger = slog.New(handler)
	})

	return s.logger
}
