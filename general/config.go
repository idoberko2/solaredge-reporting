package general

import (
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"io/fs"
	"os"
	"path"
	"runtime"

	log "github.com/sirupsen/logrus"
)

const EnvAppPrefix = "sem"

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
