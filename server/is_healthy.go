package server

import (
	"encoding/json"
	"github.com/idoberko2/semonitor/engine"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func IsHealthy(e engine.Engine) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("IsHealthy called")

		jsonResponse, err := json.Marshal(IsHealthyResponse{
			IsHealthy: e.IsHealthy(),
		})
		if err != nil {
			log.Fatal("error marshalling json")
		}

		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}
}
