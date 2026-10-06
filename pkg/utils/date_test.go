package utils

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func parseDate(t *testing.T, year, month, day int) time.Time {
	loc, err := time.LoadLocation("Asia/Singapore")
	assert.Nil(t, err)
	ret, err := time.ParseInLocation(
		time.DateTime,
		fmt.Sprintf("%d-%02d-%02d 16:00:00", year, month, day),
		loc,
	)
	assert.Nil(t, err)
	return ret
}

func TestSetDateToEndOfMonth(t *testing.T) {
	var tests = []struct {
		name     string
		given    time.Time
		expected time.Time
	}{
		{
			name:     "mid march 2023",
			given:    parseDate(t, 2023, 03, 15),
			expected: parseDate(t, 2023, 03, 31),
		},
		{
			name:     "mid Dec 2023",
			given:    parseDate(t, 2023, 12, 15),
			expected: parseDate(t, 2023, 12, 31),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, SetDateToEndOfMonth(tt.given))
		})
	}
}

func TestSetDateToDateOnly(t *testing.T) {
	singapore, err := time.LoadLocation("Asia/Singapore")
	assert.Nil(t, err)

	var tests = []struct {
		name     string
		given    time.Time
		expected string
	}{
		{
			name:     "midnight in singapore read back as utc",
			given:    time.Date(2023, 4, 3, 0, 0, 0, 0, singapore).UTC(),
			expected: "2023-04-03",
		},
		{
			name:     "evening in singapore read back as utc",
			given:    time.Date(2023, 4, 30, 20, 0, 0, 0, time.UTC),
			expected: "2023-05-01",
		},
		{
			name:     "already in singapore",
			given:    time.Date(2023, 4, 3, 12, 0, 0, 0, singapore),
			expected: "2023-04-03",
		},
		{
			name:     "already in utc during singapore business hours",
			given:    time.Date(2023, 4, 3, 4, 0, 0, 0, time.UTC),
			expected: "2023-04-03",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, SetDateToDateOnly(tt.given))
		})
	}
}
