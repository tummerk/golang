package config

import "time"

type Http struct {
	ListenAddr      string        `env:"LISTEN_ADDR" envDefault:":8080"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"10s"`
}
