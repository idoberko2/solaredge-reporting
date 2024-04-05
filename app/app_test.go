package app

import (
	"context"
	"fmt"
	"github.com/idoberko2/semonitor/general"
	"github.com/idoberko2/semonitor/server"
	"github.com/imroc/req/v3"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite
	appCancel context.CancelFunc
	srvPort   int
	c         *req.Client
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) SetupSuite() {
	general.InitBasePath()

	suite.c = req.C().SetTimeout(200 * time.Millisecond)
	suite.setupEnv()
}

func (suite *IntegrationTestSuite) SetupTest() {
	a := New()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	suite.appCancel = cancel
	go a.RunServer(ctx)
	suite.waitUntilServerReady()
}

func (suite *IntegrationTestSuite) TearDownTest() {
	suite.appCancel()
}

func (suite *IntegrationTestSuite) waitUntilServerReady() {
	ticker := time.NewTicker(100 * time.Millisecond)
	timeout := time.After(3 * time.Second)
	for {
		select {
		case <-ticker.C:
			if suite.isServerHealthy() {
				log.Info("server is ready")
				return
			}
		case <-timeout:
			suite.Assert().FailNow("waited too long for server to be ready")
		}
	}
}

func (suite *IntegrationTestSuite) TestIntegration() {
	suite.Assert().True(suite.isServerHealthy())
}

func (suite *IntegrationTestSuite) setupEnv() {
	listener, err := net.Listen("tcp", ":0")
	suite.Require().NoError(err)
	suite.srvPort = listener.Addr().(*net.TCPAddr).Port
	suite.Require().NoError(listener.Close())
	suite.Require().NoError(os.Setenv("PORT", strconv.Itoa(suite.srvPort)))
	log.WithField("port", suite.srvPort).Info("test env setup")
}

func (suite *IntegrationTestSuite) isServerHealthy() bool {
	var resp server.IsHealthyResponse

	result, err := suite.c.R().SetSuccessResult(&resp).Get(fmt.Sprintf("http://localhost:%d/api/is_healthy", suite.srvPort))
	if err != nil {
		log.WithError(err).Info("server not ready")
		return false
	}
	if result.StatusCode != http.StatusOK {
		log.WithField("status code", result.StatusCode).Info("invalid response code")
		return false
	}

	return resp.IsHealthy
}
