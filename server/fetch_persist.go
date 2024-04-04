package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/idoberko2/semonitor/engine"
	log "github.com/sirupsen/logrus"
)

const (
	errMsgEmpty         = "field '%s' cannot be empty"
	errMsgInvalidFormat = "field '%s' should be in the format 'YYYY-MM-DD'"
)

var (
	ErrFromEmpty   = fmt.Errorf(errMsgEmpty, QueryParamFrom)
	ErrToEmpty     = fmt.Errorf(errMsgEmpty, QueryParamTo)
	ErrFromInvalid = fmt.Errorf(errMsgInvalidFormat, QueryParamFrom)
	ErrToInvalid   = fmt.Errorf(errMsgInvalidFormat, QueryParamTo)
	emptyErrs      = map[string]error{
		QueryParamFrom: ErrFromEmpty,
		QueryParamTo:   ErrToEmpty,
	}
	invalidErrs = map[string]error{
		QueryParamFrom: ErrFromInvalid,
		QueryParamTo:   ErrToInvalid,
	}
)

func GetFetchPersist(e engine.Engine) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("GetFetchPersist called")
		from, to, err := parseRangeQueryParams(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(FetchPersistResponse{
				Status:  StatusError,
				Message: err.Error(),
			})
			return
		}

		if err := e.FetchAndPersist(r.Context(), from, to); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(FetchPersistResponse{
				Status:  StatusError,
				Message: "caught error during processing: " + err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(FetchPersistResponse{
			Status: StatusSuccess,
		})
	}
}

func parseRangeQueryParams(r *http.Request) (time.Time, time.Time, error) {
	from, err := parseQueryParam(r, QueryParamFrom)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	to, err := parseQueryParam(r, QueryParamTo)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return from, to, nil
}

func parseQueryParam(r *http.Request, key string) (time.Time, error) {
	rawVal := r.URL.Query().Get(key)
	if rawVal == "" {
		return time.Time{}, emptyErrs[key]
	}

	val, err := time.Parse("2006-01-02", rawVal)
	if err != nil {
		log.WithField("key", key).WithError(err).Warn("failed to parse date")
		return time.Time{}, invalidErrs[key]
	}

	return val, nil
}
