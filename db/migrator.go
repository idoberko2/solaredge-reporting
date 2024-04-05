package db

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/idoberko2/semonitor/general"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strings"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migrator interface {
	Migrate(cfg general.Config) error
}

type migrator struct {
}

func NewMigrator() Migrator {
	return &migrator{}
}

func (m *migrator) Migrate(cfg general.Config) error {
	db, err := ConnectToDb(cfg)
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	mig, err := migrate.NewWithDatabaseInstance("file://db/migrations", getDbName(cfg), driver)
	if err != nil {
		return err
	}

	log.Info("running db migration")
	migErr := mig.Up()
	if err != nil && !errors.Is(migErr, migrate.ErrNoChange) {
		return err
	}

	log.WithField("changes", !errors.Is(migErr, migrate.ErrNoChange)).Info("done migrating")

	return nil
}

func getDbName(cfg general.Config) string {
	parts := strings.Split(cfg.DbConString, "/")
	return parts[len(parts)-1]
}
