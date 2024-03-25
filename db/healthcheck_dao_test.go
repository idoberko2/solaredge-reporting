package db

import (
	"github.com/idoberko2/semonitor/general"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	general.InitBasePath()
	require.NoError(t, general.LoadDotEnv())
	cfg, err := ReadDbConfig()
	require.NoError(t, err)
	mig := NewMigrator()
	require.NoError(t, mig.Migrate(cfg))

	dao := NewHealthCheckDao(cfg)
	require.NoError(t, dao.Init())

	for i := 0; i < 3; i++ {
		assert.True(t, dao.IsHealthy())
	}
}
