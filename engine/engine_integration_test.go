//go:build integration

package engine

import (
	"context"
	"github.com/idoberko2/semonitor/db"
	"github.com/idoberko2/semonitor/general"
	"github.com/idoberko2/semonitor/seclient"
	"github.com/imroc/req/v3"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type IntegrationTestSuite struct {
	suite.Suite
	srvPort int
	dao     db.EnergyDao
	c       db.Cleaner
	svc     EnergyService
	engine  Engine
}

func TestEngineSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) SetupSuite() {
	general.InitBasePath()
	suite.Require().NoError(general.LoadDotEnv())
	dbCfg, err := db.ReadDbConfig()
	suite.Require().NoError(err)
	dao := db.NewEnergyDao(dbCfg)
	suite.Require().NoError(dao.Init())
	suite.Require().NoError(dao.Migrate())
	suite.dao = dao
	suite.c = db.NewCleaner(dbCfg)

	engineCfg, err := ReadConfig()
	suite.Require().NoError(err)

	client := req.C()
	suite.svc = NewEnergyService(
		dao,
		seclient.NewSEClient(client, engineCfg.SolarEdgeApiKey, engineCfg.SolarEdgeSiteId),
	)

	suite.engine = New(engineCfg, suite.svc)
}

func (suite *IntegrationTestSuite) SetupTest() {
	suite.Require().NoError(suite.c.Cleanup())
}

func (suite *IntegrationTestSuite) TestFetchAndPersist() {
	err := suite.engine.FetchAndPersist(context.Background(), suite.time("2024-01-01T00:00:00Z"),
		suite.time("2024-03-14T00:00:00Z"))
	suite.Require().NoError(err)

	energies, err := suite.dao.ReadEnergy(suite.time("2024-03-01T11:30:00Z"), suite.time("2024-03-01T11:31:00Z"))
	suite.Require().NoError(err)
	suite.Assert().Equal(1, len(energies))
	suite.Assert().Equal(3559, energies[0].Value)

	afterPeriodEnergies, err := suite.dao.ReadEnergy(suite.time("2024-03-15T00:00:00Z"), suite.time("2024-03-16T00:00:00Z"))
	suite.Require().NoError(err)
	suite.Assert().Len(afterPeriodEnergies, 0)
}

func (suite *IntegrationTestSuite) time(s string) time.Time {
	dt, err := time.Parse(time.RFC3339, s)
	suite.Require().NoError(err)

	return dt
}
