package engine

import (
	"context"
	"github.com/stretchr/testify/mock"
	"time"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) FetchAndPersist(ctx context.Context, from time.Time, to time.Time) error {
	args := m.Called(ctx, from, to)
	return args.Error(0)
}

func (m *Mock) IsHealthy() bool {
	args := m.Called()
	return args.Bool(0)
}
