package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"log"
)

var TIME_NEXT_TAKINGS = 1000 //в минутах

var gRPC_PORT = "12345"

var ConnStr = "postgres://postgres:pass123@localhost:3242/test?sslmode=disable"

var Key = []byte("8a1f3d9c7b2e45f60a9e8d2b4c3fds76")

type Config struct {
	Jwe      Jwe `envPrefix:"JWE_"`
	Http     Http
	Grpc     Grpc
	Service  Service
	Postgres Postgres
	Logger   Logger
}

func Load() (Config, error) {
	var config Config
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	if err = env.Parse(&config); err != nil {
		return Config{}, err
	}
	return config, nil
}
