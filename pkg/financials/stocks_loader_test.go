package financials

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func dateFromString(t *testing.T, dateString string) time.Time {
	t.Helper()
	d, err := time.Parse(time.DateOnly, dateString)
	assert.Nil(t, err)
	return d
}

func TestGetStartDate(t *testing.T) {
	earliestNeeded := dateFromString(t, "2021-03-11")

	tests := []struct {
		name            string
		earliestCovered *time.Time
		lastProcessed   *time.Time
		expected        time.Time
	}{
		{
			name:     "no coverage backfills",
			expected: earliestNeeded.AddDate(0, 0, -backfillDays),
		},
		{
			name:            "coverage starts after earliest needed backfills",
			earliestCovered: datePtr(dateFromString(t, "2021-09-24")),
			expected:        earliestNeeded.AddDate(0, 0, -backfillDays),
		},
		{
			name:            "coverage reaches earliest needed fetches incrementally",
			earliestCovered: datePtr(dateFromString(t, "2021-03-11")),
			lastProcessed:   datePtr(dateFromString(t, "2026-09-24")),
			expected:        dateFromString(t, "2026-09-17"),
		},
		{
			name:            "covered but no last processed backfills",
			earliestCovered: datePtr(dateFromString(t, "2021-03-11")),
			expected:        earliestNeeded.AddDate(0, 0, -backfillDays),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := getStartDate(tt.earliestCovered, tt.lastProcessed, earliestNeeded)
			assert.Equal(t, tt.expected, start)
		})
	}
}

func datePtr(d time.Time) *time.Time {
	return &d
}
