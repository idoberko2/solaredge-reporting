package engine

import (
	"github.com/idoberko2/semonitor/db"
	"github.com/idoberko2/semonitor/general"
	"github.com/idoberko2/semonitor/seclient"
	"github.com/imroc/req/v3"
	"github.com/jarcoal/httpmock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"net/http"
	"os"
	"testing"
	"time"
)

type EnergyServiceSuite struct {
	suite.Suite
	dao db.EnergyDao
	c   db.Cleaner
	svc EnergyService
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EnergyServiceSuite))
}

func (suite *EnergyServiceSuite) SetupSuite() {
	log.SetLevel(log.DebugLevel)
	general.InitBasePath()
	suite.Require().NoError(general.LoadDotEnv())
	cfg, err := general.ReadConfigFromEnv()
	suite.Require().NoError(err)
	mig := db.NewMigrator()
	suite.Require().NoError(mig.Migrate(cfg))

	dao := db.NewEnergyDao(cfg)
	suite.Require().NoError(dao.Init())
	suite.dao = dao
	suite.c = db.NewCleaner(cfg)
	suite.svc = NewEnergyService(dao, suite.initMockSeClient())
}

func (suite *EnergyServiceSuite) SetupTest() {
	suite.Require().NoError(suite.c.Cleanup())
}

func (suite *EnergyServiceSuite) TestRequestEnergies() {
	tf := suite.time("2024-03-01T10:00:00")
	eng, err := suite.svc.RequestEnergies(tf, tf)
	suite.Require().NoError(err)

	suite.Assert().Len(eng, 24)

	someValue := eng[22]
	expectedDt, err := time.Parse(time.RFC3339, "2024-03-01T11:45:00+02:00")
	loc, err := time.LoadLocation("Asia/Jerusalem")
	suite.Require().NoError(err)
	expectedDt = expectedDt.In(loc)
	suite.Require().NoError(err)

	suite.Assert().Equal(3513, someValue.Value)
	suite.Assert().Equal(expectedDt, someValue.DateTime)
}

func (suite *EnergyServiceSuite) TestEmptyWrites() {
	suite.Require().NoError(suite.svc.WriteEnergy([]general.Energy{}))
}

func (suite *EnergyServiceSuite) TestSimpleWrites() {
	expected := []general.Energy{{
		DateTime: suite.time("2024-03-01T12:00:00"),
		Value:    1000,
	}, {
		DateTime: suite.time("2024-03-01T12:15:00"),
		Value:    1100,
	}}

	suite.Require().NoError(suite.svc.WriteEnergy(expected))

	actual, err := suite.dao.ReadEnergy(suite.time("2024-03-01T12:00:00"), suite.time("2024-03-01T13:00:00"))
	suite.Require().NoError(err)

	suite.Assert().Equal(expected, actual)
}

func (suite *EnergyServiceSuite) TestWritesUpdates() {
	part1 := []general.Energy{{
		DateTime: suite.time("2024-03-01T10:00:00"),
		Value:    1000,
	}, {
		DateTime: suite.time("2024-03-01T10:15:00"),
		Value:    1100,
	}}

	suite.Require().NoError(suite.svc.WriteEnergy(part1))

	part2 := []general.Energy{{
		DateTime: suite.time("2024-03-01T10:15:00"),
		Value:    1120,
	}, {
		DateTime: suite.time("2024-03-01T10:30:00"),
		Value:    1200,
	}}

	suite.Require().NoError(suite.svc.WriteEnergy(part2))

	actual, err := suite.dao.ReadEnergy(suite.time("2024-03-01T10:00:00"), suite.time("2024-03-01T11:00:00"))
	suite.Require().NoError(err)

	suite.Assert().Equal([]general.Energy{part1[0], part2[0], part2[1]}, actual)
}

func (suite *EnergyServiceSuite) TestWriteMiddleUpdateNoChange() {
	part1 := []general.Energy{{
		DateTime: suite.time("2024-03-01T12:00:00"),
		Value:    1000,
	}, {
		DateTime: suite.time("2024-03-01T12:15:00"),
		Value:    1100,
	}}

	suite.Require().NoError(suite.svc.WriteEnergy(part1))

	part2 := []general.Energy{{
		DateTime: suite.time("2024-03-01T12:00:00"),
		Value:    1120,
	}, {
		DateTime: suite.time("2024-03-01T12:15:00"),
		Value:    1200,
	}}

	suite.Require().NoError(suite.svc.WriteEnergy(part2))

	actual, err := suite.dao.ReadEnergy(suite.time("2024-03-01T12:00:00"), suite.time("2024-03-01T13:00:00"))
	suite.Require().NoError(err)

	suite.Assert().Equal([]general.Energy{part1[0], part2[1]}, actual)
}

func (suite *EnergyServiceSuite) time(s string) time.Time {
	dt, err := time.Parse(time.RFC3339, s+"+02:00")
	suite.Require().NoError(err)

	loc, err := time.LoadLocation("Asia/Jerusalem")
	suite.Require().NoError(err)

	return dt.In(loc)
}

func (suite *EnergyServiceSuite) initMockSeClient() seclient.SEClient {
	general.InitBasePath()
	bytesJSON, err := os.ReadFile("./seclient/example_response_march_1st.json")
	suite.Require().NoError(err)
	respBody := string(bytesJSON)

	client := req.C()
	httpmock.ActivateNonDefault(client.GetClient())
	httpmock.RegisterResponder("GET", "https://monitoringapi.solaredge.com/site/someSiteId/energy", func(request *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(http.StatusOK, respBody)
		resp.Header.Set("Content-Type", "application/json; charset=utf-8")
		return resp, nil
	})
	return seclient.NewSEClient(client, "someApiKey", "someSiteId")
}
