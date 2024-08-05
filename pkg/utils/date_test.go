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
