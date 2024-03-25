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
	srvPort  int
	dao      db.EnergyDao
	c        db.Cleaner
	svc      EnergyService
	engine   Engine
	apiCalls []apiParamsCalls
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

	dao := db.NewEnergyDao(dbCfg)
	suite.Require().NoError(dao.Init())
	suite.dao = dao
	suite.c = db.NewCleaner(dbCfg)

	engineCfg, err := ReadConfig()
	suite.Require().NoError(err)

	suite.svc = NewEnergyService(dao, suite.initMockSeClient())
	suite.engine = New(engineCfg, suite.svc)

	suite.apiCalls = []apiParamsCalls{}
}

func (suite *MockedClientTestSuite) SetupTest() {
	suite.Require().NoError(suite.c.Cleanup())
}

func (suite *MockedClientTestSuite) TestFetchAndPersist_range() {
	err := suite.engine.FetchAndPersist(context.Background(), suite.time("2024-02-29T00:00:00Z"),
		suite.time("2024-03-01T00:00:00Z"))
	suite.Require().NoError(err)

	energies, err := suite.dao.ReadEnergy(suite.time("2024-03-01T11:30:00Z"), suite.time("2024-03-01T11:31:00Z"))
	suite.Require().NoError(err)
	suite.Assert().Equal(1, len(energies))
	suite.Assert().Equal(3559, energies[0].Value)

	afterPeriodEnergies, err := suite.dao.ReadEnergy(suite.time("2024-03-15T00:00:00Z"), suite.time("2024-03-16T00:00:00Z"))
	suite.Require().NoError(err)
	suite.Assert().Len(afterPeriodEnergies, 0)
}

func (suite *MockedClientTestSuite) TestFetchAndPersist_zeros() {
	err := suite.engine.FetchAndPersist(context.Background(), suite.time("2024-03-01T00:00:00Z"),
		suite.time("2024-03-01T00:00:00Z"))
	suite.Require().NoError(err)

	energies, err := suite.dao.ReadEnergy(suite.time("2024-03-01T23:00:00Z"), suite.time("2024-03-01T23:30:00Z"))
	suite.Require().NoError(err)
	suite.Assert().Equal(0, len(energies), "zeros should not be stored to db")
}

func (suite *MockedClientTestSuite) time(s string) time.Time {
	dt, err := time.Parse(time.RFC3339, s)
	suite.Require().NoError(err)

	return dt
}

func (suite *MockedClientTestSuite) initMockSeClient() seclient.SEClient {
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
		suite.apiCalls = append(suite.apiCalls, params)
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

func (suite *MockedClientTestSuite) readJsonStringFile(path string) string {
	bytesJSON, err := os.ReadFile(path)
	suite.Require().NoError(err)
	return string(bytesJSON)
}

type apiParamsCalls struct {
	startDate string
	endDate   string
}
