package general

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"io/fs"
	"os"
	"path"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
)

const EnvAppPrefix = "sem"

type Config struct {
	Host                  string        `default:"localhost"`
	Port                  int           `envconfig:"PORT" required:"true"`
	WriteTimeout          time.Duration `split_words:"true" default:"10s"`
	ReadTimeout           time.Duration `split_words:"true" default:"10s"`
	IdleTimeout           time.Duration `split_words:"true" default:"60s"`
	ServerShutdownTimeout time.Duration `split_words:"true" default:"10s"`
	DbConString           string        `envconfig:"DATABASE_URL"`
	SolarEdgeApiKey       string        `split_words:"true"`
	SolarEdgeSiteId       string        `split_words:"true"`
	SolarEdgeStartDate    time.Time     `split_words:"true"`
}

func ReadConfigFromEnv() (Config, error) {
	var cfg Config

	if err := envconfig.Process(EnvAppPrefix, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func InitBasePath() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func LoadDotEnv() error {
	var pathErr *fs.PathError

	if err := godotenv.Load(".env"); errors.As(err, &pathErr) {
		log.Info("couldn't find .env file, skipping .env file load")
	} else if err != nil {
		return err
	}

	return nil
}
