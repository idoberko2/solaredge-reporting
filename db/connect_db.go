package db

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func ConnectToDb(config DbConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", config.DbConString)
	if err != nil {
		return nil, err
	}

	return db, nil
}
