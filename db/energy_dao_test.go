package db

import (
	"github.com/idoberko2/semonitor/general"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type EnergyDaoSuite struct {
	suite.Suite
	dao EnergyDao
	c   Cleaner
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EnergyDaoSuite))
}

func (suite *EnergyDaoSuite) SetupSuite() {
	general.InitBasePath()
	suite.Require().NoError(general.LoadDotEnv())
	cfg, err := ReadDbConfig()
	suite.Require().NoError(err)
	mig := NewMigrator()
	suite.Require().NoError(mig.Migrate(cfg))

	dao := NewEnergyDao(cfg)
	suite.Require().NoError(dao.Init())

	suite.dao = dao
	suite.c = NewCleaner(cfg)
}

func (suite *EnergyDaoSuite) SetupTest() {
	suite.Require().NoError(suite.c.Cleanup())
}

func (suite *EnergyDaoSuite) TestReadWrite() {
	empty, err := suite.dao.ReadEnergy(suite.time("2024-03-01T10:00:00"), suite.time("2024-03-01T11:00:00"))
	suite.Require().NoError(err)
	suite.Assert().True(len(empty) == 0)

	expected := []general.Energy{{
		DateTime: suite.time("2024-03-01T10:00:00"),
		Value:    1000,
	}, {
		DateTime: suite.time("2024-03-01T10:15:00"),
		Value:    1100,
	}}
	suite.Assert().NoError(suite.dao.WriteEnergy(expected))
	actual, err := suite.dao.ReadEnergy(suite.time("2024-03-01T10:00:00"), suite.time("2024-03-01T11:00:00"))
	suite.Require().NoError(err)

	suite.Assert().Equal(expected, actual)
}

func (suite *EnergyDaoSuite) TestUpdate() {
	suite.Assert().NoError(suite.dao.WriteEnergy([]general.Energy{{
		DateTime: suite.time("2024-03-01T10:00:00"),
		Value:    1000,
	}}))

	expected := general.Energy{
		DateTime: suite.time("2024-03-01T10:00:00"),
		Value:    1100,
	}
	suite.Require().NoError(suite.dao.UpdateEnergy(expected))

	actual, err := suite.dao.ReadEnergy(suite.time("2024-03-01T10:00:00"), suite.time("2024-03-01T11:00:00"))
	suite.Require().NoError(err)

	suite.Assert().Equal([]general.Energy{expected}, actual)
}

func (suite *EnergyDaoSuite) time(s string) time.Time {
	dt, err := time.Parse(time.RFC3339, s+"+02:00")
	suite.Require().NoError(err)

	loc, err := time.LoadLocation("Asia/Jerusalem")
	suite.Require().NoError(err)

	return dt.In(loc)
}
