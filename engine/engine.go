package engine

import (
	"context"
	"github.com/idoberko2/semonitor/db"
	log "github.com/sirupsen/logrus"
	"time"
)

type Engine interface {
	FetchAndPersist(ctx context.Context, from time.Time, to time.Time) error
	IsHealthy() bool
}

func New(cfg Config, energySvc EnergyService, hcDao db.HealthCheckDao) Engine {
	return &engine{
		cfg:         cfg,
		energySvc:   energySvc,
		healthcheck: hcDao,
	}
}

type engine struct {
	cfg         Config
	energySvc   EnergyService
	healthcheck db.HealthCheckDao
	ready       bool
}

func (e *engine) FetchAndPersist(ctx context.Context, from time.Time, to time.Time) error {
	for ptr := from; !ptr.After(to); ptr = ComputeStartNextMonth(ptr) {
		end := Min(ComputeEndOfMonth(ptr), to)
		log.WithField("start", ptr).WithField("end", end).Debug("getting energy...")
		energies, err := e.energySvc.RequestEnergies(ptr, end)
		if err != nil {
			log.WithField("start", ptr).WithField("end", end).Error("failed getting energy")
			return err
		} else {
			log.
				WithField("start", ptr).
				WithField("end", end).
				WithField("energyCount", len(energies)).
				Info("got energy")
		}

		if err := e.energySvc.WriteEnergy(energies); err != nil {
			log.WithField("start", ptr).WithField("end", end).Error("failed writing energy")
			return err
		}
	}

	return nil
}

func (e *engine) IsHealthy() bool {
	return e.healthcheck.IsHealthy()
}
