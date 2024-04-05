package engine

import (
	"context"
	"github.com/idoberko2/semonitor/general"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type Engine interface {
	FetchAndPersist(ctx context.Context, from time.Time, to time.Time) error
	FetchAndPersistLastDays(ctx context.Context, rawDays string) error
}

func New(cfg general.Config, energySvc EnergyService) Engine {
	return &engine{
		cfg:       cfg,
		energySvc: energySvc,
	}
}

type engine struct {
	cfg       general.Config
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

func (e *engine) FetchAndPersistLastDays(ctx context.Context, rawDays string) error {
	days, err := strconv.Atoi(rawDays)
	if err != nil {
		return err
	}

	to := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
	from := to.Add(-24 * time.Duration(days) * time.Hour)

	return e.FetchAndPersist(ctx, from, to)
}
