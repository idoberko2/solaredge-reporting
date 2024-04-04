//go:build integration

package engine

import (
	"context"
	"github.com/idoberko2/semonitor/db"
	"github.com/idoberko2/semonitor/general"
	"github.com/idoberko2/semonitor/seclient"
	"github.com/imroc/req/v3"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"os"
	"testing"
	"time"
)

type MockedClientTestSuite struct {
	suite.Suite
	srvPort   int
	dao       db.EnergyDao
	c         db.Cleaner
	engineCfg Config
}

func TestEngineMockedClientSuite(t *testing.T) {
	suite.Run(t, new(MockedClientTestSuite))
}

func (suite *MockedClientTestSuite) SetupSuite() {
	general.InitBasePath()
	suite.Require().NoError(general.LoadDotEnv())
	dbCfg, err := db.ReadDbConfig()
	suite.Require().NoError(err)
	mig := db.NewMigrator()
	suite.Require().NoError(mig.Migrate(dbCfg))

	eDao := db.NewEnergyDao(dbCfg)
	suite.Require().NoError(eDao.Init())
	suite.dao = eDao
	suite.c = db.NewCleaner(dbCfg)

	engineCfg, err := ReadConfig()
	suite.Require().NoError(err)
	suite.engineCfg = engineCfg
}

func (suite *MockedClientTestSuite) SetupTest() {
	suite.Require().NoError(suite.c.Cleanup())
}

func (suite *MockedClientTestSuite) TestFetchAndPersist_range() {
	engine := New(suite.engineCfg, NewEnergyService(suite.dao, suite.initMockSeSeparateDaysClient()))
	err := engine.FetchAndPersist(
		context.Background(),
		time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 3, 01, 0, 0, 0, 0, time.UTC),
	)
	suite.Require().NoError(err)

	febEnergies, err := suite.dao.ReadEnergy(suite.time("2024-02-29T14:00:00"), suite.time("2024-02-29T14:01:00"))
	suite.Require().NoError(err)
	suite.Require().Equal(1, len(febEnergies))
	suite.Assert().Equal(3334, febEnergies[0].Value)

	marchEnergies, err := suite.dao.ReadEnergy(suite.time("2024-03-01T11:30:00"), suite.time("2024-03-01T11:31:00"))
	suite.Require().NoError(err)
	suite.Require().Equal(1, len(marchEnergies))
	suite.Assert().Equal(3559, marchEnergies[0].Value)

	afterPeriodEnergies, err := suite.dao.ReadEnergy(suite.time("2024-03-15T00:00:00"), suite.time("2024-03-16T00:00:00"))
	suite.Require().NoError(err)
	suite.Assert().Len(afterPeriodEnergies, 0)
}

func (suite *MockedClientTestSuite) TestFetchAndPersist_zeros() {
	engine := New(suite.engineCfg, NewEnergyService(suite.dao, suite.initMockSeSeparateDaysClient()))
	err := engine.FetchAndPersist(
		context.Background(),
		time.Date(2024, 3, 01, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 3, 01, 0, 0, 0, 0, time.UTC),
	)
	suite.Require().NoError(err)

	energies, err := suite.dao.ReadEnergy(suite.time("2024-03-01T23:00:00"), suite.time("2024-03-01T23:30:00"))
	suite.Require().NoError(err)
	suite.Assert().Equal(0, len(energies), "zeros should not be stored to db")
}

func (suite *MockedClientTestSuite) TestFetchAndPersist_update() {
	engine := New(suite.engineCfg, NewEnergyService(suite.dao, suite.initMockSeSameDayUpdateClient()))
	err := engine.FetchAndPersist(
		context.Background(),
		time.Date(2024, 3, 01, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 3, 01, 0, 0, 0, 0, time.UTC),
	)
	suite.Require().NoError(err)

	// "later" response with an update and more entries
	suite.Require().NoError(engine.FetchAndPersist(
		context.Background(),
		time.Date(2024, 3, 01, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 3, 01, 0, 0, 0, 0, time.UTC),
	))

	energies, err := suite.dao.ReadEnergy(suite.time("2024-03-01T11:00:00"), suite.time("2024-03-01T12:01:00"))
	suite.Require().NoError(err)
	suite.Require().Equal(5, len(energies))
	suite.Assert().Equal(3344, energies[0].Value, "expected value to be updated")
	suite.Assert().Equal(3328, energies[1].Value, "expected new entry to be stored")
	suite.Assert().Equal(3559, energies[2].Value, "expected new entry to be stored")
	suite.Assert().Equal(3513, energies[3].Value, "expected new entry to be stored")
	suite.Assert().Equal(2233, energies[4].Value, "expected new entry to be stored")
}

func (suite *MockedClientTestSuite) time(s string) time.Time {
	dt, err := time.Parse(time.RFC3339, s+"+02:00")
	suite.Require().NoError(err)

	loc, err := time.LoadLocation("Asia/Jerusalem")
	suite.Require().NoError(err)

	return dt.In(loc)
}

func (suite *MockedClientTestSuite) initMockSeSeparateDaysClient() seclient.SEClient {
	feb29thRespBody := suite.readJsonStringFile("./seclient/example_response_feb_29th.json")
	march1stRespBody := suite.readJsonStringFile("./seclient/example_response_march_1st.json")

	client := req.C()
	httpmock.ActivateNonDefault(client.GetClient())
	httpmock.RegisterResponder("GET", "https://monitoringapi.solaredge.com/site/someSiteId/energy", func(request *http.Request) (*http.Response, error) {
		var body string

		params := apiParamsCalls{
			startDate: request.URL.Query().Get("startDate"),
			endDate:   request.URL.Query().Get("endDate"),
		}
		if params.startDate == "2024-02-29" && params.endDate == "2024-02-29" {
			body = feb29thRespBody
		} else if params.startDate == "2024-03-01" && params.endDate == "2024-03-01" {
			body = march1stRespBody
		} else {
			suite.FailNowf("unexpected dates requested from mock", "startDate=%s, endDate=%s",
				params.startDate, params.endDate)
		}

		resp := httpmock.NewStringResponse(http.StatusOK, body)
		resp.Header.Set("Content-Type", "application/json; charset=utf-8")
		return resp, nil
	})

	return seclient.NewSEClient(client, "someApiKey", "someSiteId")
}

func (suite *MockedClientTestSuite) initMockSeSameDayUpdateClient() seclient.SEClient {
	earlyRespBody := suite.readJsonStringFile("./seclient/example_response_march_early.json")
	lateRespBody := suite.readJsonStringFile("./seclient/example_response_march_late.json")

	client := req.C()
	httpmock.ActivateNonDefault(client.GetClient())

	responses := []*http.Response{
		httpmock.NewStringResponse(http.StatusOK, earlyRespBody),
		httpmock.NewStringResponse(http.StatusOK, lateRespBody),
	}

	for _, resp := range responses {
		resp.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	httpmock.ResponderFromMultipleResponses(responses)
	httpmock.RegisterResponder(
		"GET",
		"https://monitoringapi.solaredge.com/site/someSiteId/energy",
		httpmock.ResponderFromMultipleResponses(responses),
	)

	return seclient.NewSEClient(client, "someApiKey", "someSiteId")
}

func (suite *MockedClientTestSuite) readJsonStringFile(path string) string {
	bytesJSON, err := os.ReadFile(path)
	suite.Require().NoError(err)
	return string(bytesJSON)
}

type apiParamsCalls struct {
	startDate string
	endDate   string
}
