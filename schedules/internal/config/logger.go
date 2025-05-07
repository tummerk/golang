package config

type Logger struct {
	Filename string `env:"LOG_FILENAME"`
	Debug    bool   `env:"LOG_DEBUG"`
}
