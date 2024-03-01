//go:build integration
// +build integration

package app

import (
	"context"
	"fmt"
	"github.com/idoberko2/semonitor/general"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

/**
This is a full E2E integration tests. It sets up the environment, sends a ping which updates the
state to "Healthy", then waits 6 seconds for the state to become "Unhealthy".

You should monitor the configured TG group when running the test. You expect to see:
1. A "Healthy" message arrives immediately
2. An "Unhealthy" message arrives 5 seconds later
*/

type IntegrationTestSuite struct {
	suite.Suite
	appCancel context.CancelFunc
	srvPort   int
	c         http.Client
}

func (suite *IntegrationTestSuite) SetupSuite() {
	suite.Require().NoError(general.LoadDotEnv())

	suite.c = http.Client{}
	suite.c.Timeout = 200 * time.Millisecond
	suite.setupEnv()
}

func (suite *IntegrationTestSuite) SetupTest() {
	a := New()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	suite.appCancel = cancel
	go a.Run(ctx)
	suite.waitUntilServerReady()
}

func (suite *IntegrationTestSuite) TearDownTest() {
	suite.appCancel()
}

func (suite *IntegrationTestSuite) waitUntilServerReady() {
	var lastErr error
	var lastStatusCode int
	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			{
				resp, err := suite.c.Get(fmt.Sprintf("http://localhost:%d/is_alive", suite.srvPort))
				lastErr = err
				if resp != nil {
					lastStatusCode = resp.StatusCode
				}

				if lastErr == nil && lastStatusCode == http.StatusOK {
					return
				}
			}
		case <-time.After(3 * time.Second):
			{
				var errStr string
				if lastErr != nil {
					errStr = lastErr.Error()
				}
				suite.Assert().FailNowf("waited too long for server to be ready err=%s status_code=%d", errStr, lastStatusCode)
			}
		}
	}
}

func (suite *IntegrationTestSuite) TestIntegration() {
	passphrase := os.Getenv("HC_PASSPHRASE")
	resp, err := suite.c.Post(fmt.Sprintf("http://localhost:%d/ping", suite.srvPort), "application/json", strings.NewReader(fmt.Sprintf("{\"passphrase\": \"%s\"}", passphrase)))
	suite.Require().NoError(err)
	suite.Assert().Equal(http.StatusOK, resp.StatusCode)

	<-time.After(6 * time.Second)
}

func (suite *IntegrationTestSuite) setupEnv() {
	listener, err := net.Listen("tcp", ":0")
	suite.Require().NoError(err)
	suite.srvPort = listener.Addr().(*net.TCPAddr).Port
	suite.Require().NoError(listener.Close())
	os.Setenv("PORT", strconv.Itoa(suite.srvPort))
	os.Setenv("HC_SAMPLE_RATE", "100ms")
	os.Setenv("HC_GRACE_PERIOD", "5s")
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
