package config

type Grpc struct {
	Addr string `env:"Addr" envDefault:":12345"`
}
