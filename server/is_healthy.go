package server

import (
	"encoding/json"
	"github.com/idoberko2/semonitor/db"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func IsHealthy(hcDao db.HealthCheckDao) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("IsHealthy called")

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(IsHealthyResponse{
			IsHealthy: hcDao.IsHealthy(),
		})
	}
}
