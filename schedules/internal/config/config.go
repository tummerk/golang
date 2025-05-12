package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"log"
)

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
