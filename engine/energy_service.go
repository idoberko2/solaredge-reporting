package engine

import (
	"github.com/idoberko2/semonitor/db"
	"github.com/idoberko2/semonitor/general"
	"github.com/idoberko2/semonitor/seclient"
	log "github.com/sirupsen/logrus"
	"time"
)

type EnergyService interface {
	RequestEnergies(start time.Time, end time.Time) ([]general.Energy, error)
	WriteEnergy([]general.Energy) error
}

func NewEnergyService(dao db.EnergyDao, secli seclient.SEClient) EnergyService {
	return &energyService{
		dao:   dao,
		secli: secli,
	}
}

type energyService struct {
	dao   db.EnergyDao
	secli seclient.SEClient
}

func (e *energyService) RequestEnergies(start time.Time, end time.Time) ([]general.Energy, error) {
	allEnergies, err := e.secli.GetEnergy(start, end)
	if err != nil {
		return nil, err
	}

	return filterNonZeroEnergies(allEnergies), nil
}

func (e *energyService) WriteEnergy(energies []general.Energy) error {
	var toWrite []general.Energy
	toUpdate := general.Energy{}
	lastPersisted := general.Energy{}

	persisted, err := e.dao.ReadEnergy(findMinDt(energies), findMaxDt(energies).Add(time.Hour))
	if err != nil {
		return err
	}

	if len(persisted) > 0 {
		lastPersisted = persisted[len(persisted)-1]
	}

	persistedValueByTime := toMapByTimestamp(persisted)

	for _, energy := range energies {
		persistedEnergy, ok := persistedValueByTime[energy.DateTime.Unix()]
		if !ok {
			log.
				WithField("datetime", energy.DateTime).
				WithField("value", energy.Value).
				Debug("adding entry to write")
			toWrite = append(toWrite, energy)
		} else if energy.DateTime.Equal(lastPersisted.DateTime) && persistedEnergy != energy.Value {
			// only last entry might be updated
			log.
				WithField("datetime", energy.DateTime).
				WithField("persistedValue", persistedEnergy).
				WithField("newValue", energy.Value).
				Info("found change in last entry value")
			toUpdate = energy
		} else if persistedEnergy != energy.Value {
			log.
				WithField("datetime", energy.DateTime).
				WithField("persistedValue", persistedEnergy).
				WithField("newValue", energy.Value).
				Error("unexpected energy change")
		}
	}

	if len(toWrite) == 0 {
		log.Info("no insert entries to write")
	} else if err := e.dao.WriteEnergy(toWrite); err != nil {
		return err
	} else {
		log.WithField("entriesCount", len(toWrite)).Info("successfully written entries")
	}

	if toUpdate.Empty() {
		log.Debug("no change in last entry, nothing to update")
	} else {
		log.
			WithField("datetime", toUpdate.DateTime).
			WithField("value", toUpdate.Value).
			Info("updated last entry")
		return e.dao.UpdateEnergy(toUpdate)
	}

	return nil
}

func filterNonZeroEnergies(allEnergies []general.Energy) []general.Energy {
	nonZero := 0

	for _, energy := range allEnergies {
		if energy.Value > 0 {
			nonZero += 1
		}
	}

	energies := make([]general.Energy, 0, nonZero)
	for _, energy := range allEnergies {
		if energy.Value > 0 {
			energies = append(energies, energy)
		} else {
			log.WithField("datetime", energy.DateTime).Debug("filtering out zero value energy")
		}
	}

	return energies
}

func toMapByTimestamp(energies []general.Energy) map[int64]int {
	res := map[int64]int{}
	for _, energy := range energies {
		res[energy.DateTime.Unix()] = energy.Value
	}

	return res
}

func findMinDt(energies []general.Energy) time.Time {
	minDt := time.Now()
	for _, energy := range energies {
		if energy.DateTime.Before(minDt) {
			minDt = energy.DateTime
		}
	}

	return minDt
}

func findMaxDt(energies []general.Energy) time.Time {
	maxDt := time.UnixMilli(0)
	for _, energy := range energies {
		if energy.DateTime.After(maxDt) {
			maxDt = energy.DateTime
		}
	}

	return maxDt
}
