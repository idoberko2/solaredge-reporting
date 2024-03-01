package db

import "github.com/jmoiron/sqlx"

type Cleaner interface {
	Cleanup() error
}

func NewCleaner(config DbConfig) Cleaner {
	db, err := ConnectToDb(config)
	if err != nil {
		panic(err)
	}

	return &cleaner{db: db}
}

type cleaner struct {
	db *sqlx.DB
}

func (c *cleaner) Cleanup() error {
	_, err := c.db.Exec(
		"TRUNCATE TABLE se_data;")
	return err
}
