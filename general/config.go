package general

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"os"
	"path"
	"runtime"
)

const EnvAppPrefix = "sem"

type AppConfig struct {
	AvoidDotEnv bool `split_words:"true"`
}

func ReadAppConfig() (AppConfig, error) {
	var cfg AppConfig

	if err := envconfig.Process(EnvAppPrefix, &cfg); err != nil {
		return cfg, errors.Wrap(err, "error processing app config")
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
	appConfig, err := ReadAppConfig()
	if err != nil {
		return err
	}

	if !appConfig.AvoidDotEnv {
		if err := godotenv.Load(); err != nil {
			return err
		}
	}

	return nil
}
