package server

import (
	"fmt"
	"github.com/idoberko2/semonitor/db"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/idoberko2/semonitor/engine"
)

func New(e engine.Engine, hcDao db.HealthCheckDao, cfg ServerConfig) *http.Server {
	r := mux.NewRouter()

	apir := r.PathPrefix("/api").Subrouter()
	apir.HandleFunc("/fetch_persist", GetFetchPersist(e)).Queries("from", "{from}", "to", "{to}").Methods(http.MethodGet)
	apir.HandleFunc("/is_healthy", IsHealthy(hcDao)).Methods(http.MethodGet)

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		Handler:      r,
	}

	return srv
}
