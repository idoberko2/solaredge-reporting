package db

import (
	"github.com/idoberko2/semonitor/general"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var ErrNotInitialized = errors.New("DB not initialized")

type EnergyDao interface {
	Init() error
	WriteEnergy(energy []general.Energy) error
	UpdateEnergy(energy general.Energy) error
	ReadEnergy(from time.Time, to time.Time) ([]general.Energy, error)
}

func NewEnergyDao(config general.Config) EnergyDao {
	return &energyDao{config: config}
}

type energyDao struct {
	config general.Config
	db     *sqlx.DB
}

func (e *energyDao) Init() error {
	db, err := ConnectToDb(e.config)
	if err != nil {
		return err
	}
	e.db = db

	return nil
}

func (e *energyDao) WriteEnergy(energy []general.Energy) error {
	if err := e.checkIsInitialized(); err != nil {
		return err
	}

	tx, err := e.db.Begin()
	if err != nil {
		return errors.Wrap(err, "error in begin transaction")
	}

	stmt, err := tx.Prepare("INSERT INTO se_data(t, value) VALUES ($1, $2);")
	if err != nil {
		return errors.Wrap(err, "error in query preparation")
	}

	for _, e := range energy {
		if _, err := stmt.Exec(e.DateTime, e.Value); err != nil {
			return errors.Wrap(err, "error in execute statement")
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "error in commit transaction")
	}
	if err := stmt.Close(); err != nil {
		return errors.Wrap(err, "error in close statement")
	}

	return nil
}

func (e *energyDao) UpdateEnergy(energy general.Energy) error {
	if err := e.checkIsInitialized(); err != nil {
		return err
	}

	_, err := e.db.Exec("UPDATE se_data SET value=$1 WHERE t=$2;", energy.Value, energy.DateTime)
	return err
}

func (e *energyDao) ReadEnergy(from time.Time, to time.Time) ([]general.Energy, error) {
	if err := e.checkIsInitialized(); err != nil {
		return nil, err
	}

	var res []general.Energy
	query := "SELECT t, value FROM se_data WHERE t >= $1 AND t < $2 ORDER BY t;"
	if err := e.db.Select(&res, query, from, to); err != nil {
		return nil, err
	}

	var finalRes []general.Energy
	loc, err := time.LoadLocation("Asia/Jerusalem")
	if err != nil {
		return nil, err
	}

	for _, entry := range res {
		finalRes = append(finalRes, general.Energy{
			DateTime: entry.DateTime.In(loc),
			Value:    entry.Value,
		})
	}

	return finalRes, nil
}

func (e *energyDao) checkIsInitialized() error {
	if !e.isInitialized() {
		return ErrNotInitialized
	}

	return nil
}

func (e *energyDao) isInitialized() bool {
	return e.db != nil
}
