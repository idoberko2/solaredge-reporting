package db

import (
	"github.com/idoberko2/semonitor/general"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type DbConfig struct {
	DbConString string `envconfig:"DATABASE_URL"`
}

func ReadDbConfig() (DbConfig, error) {
	var cfg DbConfig

	if err := envconfig.Process(general.EnvAppPrefix, &cfg); err != nil {
		return cfg, errors.Wrap(err, "error processing db config")
	}

	return cfg, nil
}
