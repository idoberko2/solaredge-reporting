package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/idoberko2/semonitor/engine"
	log "github.com/sirupsen/logrus"
)

func GetFetchPersist(e engine.Engine) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("GetFetchPersist called")
		fromRaw := r.URL.Query().Get("from")
		if fromRaw == "" {
			jsonResponse, err := json.Marshal(FetchPersistResponse{
				Status:  StatusError,
				Message: "Field 'from' cannot be empty",
			})
			if err != nil {
				log.Fatal("error marshalling json")
			}

			w.WriteHeader(http.StatusBadRequest)
			w.Write(jsonResponse)
			return
		}
		from, err := time.Parse("2006-01-02", fromRaw)
		if err != nil {
			jsonResponse, err := json.Marshal(FetchPersistResponse{
				Status:  StatusError,
				Message: "Field 'from' should be in the format 'YYYY-MM-DD'",
			})
			if err != nil {
				log.Fatal("error marshalling json")
			}

			w.WriteHeader(http.StatusBadRequest)
			w.Write(jsonResponse)
			return
		}

		toRaw := r.URL.Query().Get("from")
		if toRaw == "" {
			jsonResponse, err := json.Marshal(FetchPersistResponse{
				Status:  StatusError,
				Message: "Field 'to' cannot be empty",
			})
			if err != nil {
				log.Fatal("error marshalling json")
			}

			w.WriteHeader(http.StatusBadRequest)
			w.Write(jsonResponse)
			return
		}
		to, err := time.Parse("2006-01-02", toRaw)
		if err != nil {
			jsonResponse, err := json.Marshal(FetchPersistResponse{
				Status:  StatusError,
				Message: "Field 'to' should be in the format 'YYYY-MM-DD'",
			})
			if err != nil {
				log.Fatal("error marshalling json")
			}

			w.WriteHeader(http.StatusBadRequest)
			w.Write(jsonResponse)
			return
		}

		if err := e.FetchAndPersist(r.Context(), from, to); err != nil {
			jsonResponse, err := json.Marshal(FetchPersistResponse{
				Status:  StatusError,
				Message: "Caught error during processing: " + err.Error(),
			})
			if err != nil {
				log.Fatal("error marshalling json")
			}

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(jsonResponse)
			return
		}

		jsonResponse, err := json.Marshal(FetchPersistResponse{
			Status: StatusSuccess,
		})
		if err != nil {
			log.Fatal("error marshalling json")
		}

		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}
}
