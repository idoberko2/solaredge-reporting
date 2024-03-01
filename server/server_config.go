package server

import "time"

type ServerConfig struct {
	Host                  string        `default:"localhost"`
	Port                  int           `envconfig:"PORT" required:"true"`
	WriteTimeout          time.Duration `split_words:"true" default:"10s"`
	ReadTimeout           time.Duration `split_words:"true" default:"10s"`
	IdleTimeout           time.Duration `split_words:"true" default:"60s"`
	ServerShutdownTimeout time.Duration `split_words:"true" default:"10s"`
}
