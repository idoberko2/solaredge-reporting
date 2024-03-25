package db

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/pkg/errors"
	"strings"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migrator interface {
	Migrate(cfg DbConfig) error
}

type migrator struct {
}

func NewMigrator() Migrator {
	return &migrator{}
}

func (m *migrator) Migrate(cfg DbConfig) error {
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

	if err := mig.Up(); !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func getDbName(cfg DbConfig) string {
	parts := strings.Split(cfg.DbConString, "/")
	return parts[len(parts)-1]
}
