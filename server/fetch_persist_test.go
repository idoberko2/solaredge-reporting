package server

import (
	"encoding/json"
	"github.com/idoberko2/semonitor/engine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestParamsParse(t *testing.T) {
	m := engine.Mock{}
	m.On("FetchAndPersist", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	h := GetFetchPersist(&m)
	ts := httptest.NewServer(http.HandlerFunc(h))
	defer ts.Close()

	_, err := http.Get(ts.URL + "?from=2024-01-14&to=2024-03-02")
	require.NoError(t, err)

	m.AssertCalled(t, "FetchAndPersist", mock.Anything,
		time.Date(2024, 1, 14, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 3, 2, 0, 0, 0, 0, time.UTC))
}

func TestParseInvalidParams(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected error
	}{
		{"empty from", "?to=2024-03-02", ErrFromEmpty},
		{"invalid from", "?from=2024-01-114&to=2024-03-02", ErrFromInvalid},
		{"empty to", "?from=2024-03-02", ErrToEmpty},
		{"invalid to", "?from=2024-01-14&to=2024-013-02", ErrToInvalid},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := engine.Mock{}
			h := GetFetchPersist(&m)
			ts := httptest.NewServer(http.HandlerFunc(h))
			defer ts.Close()

			var respPayload FetchPersistResponse
			resp, err := http.Get(ts.URL + test.query)
			require.NoError(t, err)
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(resp.Body)

			require.NoError(t, json.NewDecoder(resp.Body).Decode(&respPayload))
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			assert.Equal(t, StatusError, respPayload.Status)
			assert.Equal(t, test.expected.Error(), respPayload.Message)
		})
	}
}
