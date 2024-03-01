package app

import (
	"github.com/idoberko2/semonitor/general"
	"github.com/idoberko2/semonitor/notifier"
	"github.com/idoberko2/semonitor/server"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

func ReadServerConfig() (server.ServerConfig, error) {
	var cfg server.ServerConfig

	if err := envconfig.Process(general.EnvAppPrefix, &cfg); err != nil {
		return cfg, errors.Wrap(err, "error processing server config")
	}

	return cfg, nil
}

func ReadTelegramConfig() (notifier.TelegramConfig, error) {
	var cfg notifier.TelegramConfig

	if err := envconfig.Process(general.EnvAppPrefix, &cfg); err != nil {
		return cfg, errors.Wrap(err, "error processing telegram config")
	}

	return cfg, nil
}

var ErrEmptyPassphrase = errors.New("passphrase is not set")
var ErrEmptyPort = errors.New("port is not set")
