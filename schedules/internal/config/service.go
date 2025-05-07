package config

type Service struct {
	TimeNextTakings int `env:"TIME_NEXT_TAKINGS" envDefault:"60"`
}
