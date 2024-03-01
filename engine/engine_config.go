package engine

import (
	"github.com/idoberko2/semonitor/general"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"time"
)

type Config struct {
	errReporter        chan error
	SolarEdgeApiKey    string    `split_words:"true"`
	SolarEdgeSiteId    string    `split_words:"true"`
	SolarEdgeStartDate time.Time `split_words:"true"`
}

func ReadConfig() (Config, error) {
	var cfg Config

	if err := envconfig.Process(general.EnvAppPrefix, &cfg); err != nil {
		return cfg, errors.Wrap(err, "error processing engine config")
	}

	return cfg, nil
}
