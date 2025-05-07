package config

type Jwe struct {
	Key string `env:"KEY,required"`
}
