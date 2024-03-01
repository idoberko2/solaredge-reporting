package engine

import (
	"context"
	log "github.com/sirupsen/logrus"
	"time"
)

type Engine interface {
	FetchAndPersist(ctx context.Context, from time.Time, to time.Time) error
}

func New(cfg Config, energySvc EnergyService) Engine {
	return &engine{
		cfg:       cfg,
		energySvc: energySvc,
	}
}

type engine struct {
	cfg       Config
	energySvc EnergyService
	ready     bool
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
