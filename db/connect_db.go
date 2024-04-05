package db

import (
	"github.com/idoberko2/semonitor/general"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func ConnectToDb(config general.Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", config.DbConString)
	if err != nil {
		return nil, err
	}

	return db, nil
}
