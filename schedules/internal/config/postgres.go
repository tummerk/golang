package config

type Postgres struct {
	DSN string `env:"DB_URL"`
}
