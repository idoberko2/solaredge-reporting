package seclient

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/imroc/req/v3"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	bytesJSON, err := os.ReadFile("./example_response_march_1st.json")
	require.NoError(t, err)
	respBody := string(bytesJSON)

	client := req.C()
	httpmock.ActivateNonDefault(client.GetClient())
	httpmock.RegisterResponder("GET", "https://monitoringapi.solaredge.com/site/someSiteId/energy", func(request *http.Request) (*http.Response, error) {
		assert.Equal(t, "QUARTER_OF_AN_HOUR", request.URL.Query().Get("timeUnit"))
		assert.Equal(t, "someApiKey", request.URL.Query().Get("api_key"))
		resp := httpmock.NewStringResponse(http.StatusOK, respBody)
		resp.Header.Set("Content-Type", "application/json; charset=utf-8")
		return resp, nil
	})
	seclient := NewSEClient(client, "someApiKey", "someSiteId")

	tf, err := time.Parse("2006-01-02", "2024-03-01")
	require.NoError(t, err)

	eng, err := seclient.GetEnergy(tf, tf)
	if err != nil {
		t.Error(err)
	}

	assert.Len(t, eng, 24*4)

	someValue := eng[47]
	expectedDt, err := time.Parse(time.RFC3339, "2024-03-01T09:45:00Z")
	loc, err := time.LoadLocation("Asia/Jerusalem")
	require.NoError(t, err)
	expectedDt = expectedDt.In(loc)
	require.NoError(t, err)

	assert.Equal(t, 3513, someValue.Value)
	assert.Equal(t, expectedDt, someValue.DateTime)
}
