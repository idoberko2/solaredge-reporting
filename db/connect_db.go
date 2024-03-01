package db

import (
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func ConnectToDb(config DbConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.DbUser, config.DbPassword, config.DbHost, config.DbPort, config.DbName))
	if err != nil {
		return nil, err
	}

	return db, nil
}
