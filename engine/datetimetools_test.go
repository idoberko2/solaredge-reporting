package engine

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestComputeStartNextMonth(t *testing.T) {
	var testCases = []struct {
		input    time.Time
		expected time.Time
	}{
		{input: createDateTime(t, "2023-10-17T00:00:00Z"), expected: createDateTime(t, "2023-11-01T00:00:00Z")},
		{input: createDateTime(t, "2023-10-17T12:31:00Z"), expected: createDateTime(t, "2023-11-01T00:00:00Z")},
		{input: createDateTime(t, "2023-12-21T00:00:00Z"), expected: createDateTime(t, "2024-01-01T00:00:00Z")},
		{input: createDateTime(t, "2024-02-29T00:00:00Z"), expected: createDateTime(t, "2024-03-01T00:00:00Z")},
	}

	for _, tt := range testCases {
		actual := ComputeStartNextMonth(tt.input)
		assert.Equal(t, tt.expected, actual, "for input "+tt.input.Format(time.RFC3339))
	}
}

func TestComputeEndOFMonth(t *testing.T) {
	var testCases = []struct {
		input    time.Time
		expected time.Time
	}{
		{input: createDateTime(t, "2023-10-17T00:00:00Z"), expected: createDateTime(t, "2023-10-31T00:00:00Z")},
		{input: createDateTime(t, "2023-10-17T12:31:00Z"), expected: createDateTime(t, "2023-10-31T00:00:00Z")},
		{input: createDateTime(t, "2023-12-21T00:00:00Z"), expected: createDateTime(t, "2023-12-31T00:00:00Z")},
		{input: createDateTime(t, "2023-12-31T00:00:00Z"), expected: createDateTime(t, "2023-12-31T00:00:00Z")},
		{input: createDateTime(t, "2024-02-25T00:00:00Z"), expected: createDateTime(t, "2024-02-29T00:00:00Z")},
	}

	for _, tt := range testCases {
		actual := ComputeEndOfMonth(tt.input)
		assert.Equal(t, tt.expected, actual, "for input "+tt.input.Format(time.RFC3339))
	}
}

func createDateTime(t *testing.T, s string) time.Time {
	dt, err := time.Parse(time.RFC3339, s)
	require.NoError(t, err)

	return dt
}
