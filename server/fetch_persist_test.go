package server

import (
	"github.com/idoberko2/semonitor/engine"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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
