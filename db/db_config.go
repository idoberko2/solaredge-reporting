package db

import (
	"github.com/idoberko2/semonitor/general"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type DbConfig struct {
	DbName     string `split_words:"true"`
	DbHost     string `split_words:"true"`
	DbPort     string `split_words:"true"`
	DbUser     string `split_words:"true"`
	DbPassword string `split_words:"true"`
}

func ReadDbConfig() (DbConfig, error) {
	var cfg DbConfig

	if err := envconfig.Process(general.EnvAppPrefix, &cfg); err != nil {
		return cfg, errors.Wrap(err, "error processing db config")
	}

	return cfg, nil
}
