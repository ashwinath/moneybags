package utils

import (
	"fmt"
	"time"
)

var financialsLocation = time.UTC

func SetDateToEndOfMonth(d time.Time) time.Time {
	year, month, _ := d.Date()

	loc, _ := time.LoadLocation("Asia/Singapore")
	ret, _ := time.ParseInLocation(
		time.DateTime,
		fmt.Sprintf("%d-%02d-01 16:00:00", year, month),
		loc,
	)
	return ret.AddDate(0, 1, 0).AddDate(0, 0, -1)
}

func SetDateToEndOfMonthFinancials(d time.Time) time.Time {
	year, month, _ := d.Date()

	ret, _ := time.ParseInLocation(
		time.DateTime,
		fmt.Sprintf("%d-%02d-01 08:00:00", year, month),
		financialsLocation,
	)
	return ret.AddDate(0, 1, 0).AddDate(0, 0, -1)
}

func GetLastDateOfMonth(d time.Time) time.Time {
	d = time.Date(d.Year(), d.Month(), 1, 8, 0, 0, 0, financialsLocation)
	d = d.AddDate(0, 1, 0)
	d = d.AddDate(0, 0, -1)

	return d
}

func GetFirstDateOfMonth(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), 1, 8, 0, 0, 0, financialsLocation)
}

// Accepts yyyy-mm-dd
func SetDateFromString(date string) (time.Time, error) {
	dateString := fmt.Sprintf("%s 08:00:00", date)
	return time.ParseInLocation(time.DateTime, dateString, financialsLocation)
}

// SG time
func SetDateTo0000Hours(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 8, 0, 0, 0, financialsLocation)
}

func CopyTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), d.Hour(), d.Minute(), d.Second(), d.Nanosecond(), financialsLocation)
}
