package db

import (
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"time"
)

type HealthCheckDao interface {
	Init() error
	IsHealthy() bool
}

type healthCheckDao struct {
	config DbConfig
	db     *sqlx.DB
}

func (h *healthCheckDao) Init() error {
	db, err := ConnectToDb(h.config)
	if err != nil {
		return err
	}
	h.db = db

	return nil
}

func (h *healthCheckDao) IsHealthy() bool {
	if err := h.checkIsInitialized(); err != nil {
		log.WithError(err).Error("healthcheck db not initialized")
		return false
	}

	if _, err := h.db.Exec("TRUNCATE healthcheck;"); err != nil {
		log.WithError(err).Error("error truncating the healthcheck table")
		return false
	}

	now := time.Now()
	if _, err := h.db.Exec("INSERT INTO healthcheck(t) VALUES ($1);", now); err != nil {
		log.WithError(err).Error("error writing to the healthcheck table")
		return false
	}

	var actual time.Time
	if err := h.db.Get(&actual, "SELECT t FROM healthcheck;"); err != nil {
		log.WithError(err).Error("error selecting from the healthcheck table")
		return false
	}

	return true
}

func NewHealthCheckDao(config DbConfig) HealthCheckDao {
	return &healthCheckDao{config: config}
}

func (h *healthCheckDao) checkIsInitialized() error {
	if !h.isInitialized() {
		return ErrNotInitialized
	}

	return nil
}

func (h *healthCheckDao) isInitialized() bool {
	return h.db != nil
}
