package utils

import (
	"fmt"
	"time"
)

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

	loc, _ := time.LoadLocation("Asia/Singapore")
	ret, _ := time.ParseInLocation(
		time.DateTime,
		fmt.Sprintf("%d-%02d-01 00:00:00", year, month),
		loc,
	)
	return ret.AddDate(0, 1, 0).AddDate(0, 0, -1)
}

func GetLastDateOfMonth(d time.Time) time.Time {
	d = time.Date(d.Year(), d.Month(), 1, d.Hour(), 0, 0, 0, d.Location())
	d = d.AddDate(0, 1, 0)
	d = d.AddDate(0, 0, -1)

	return d
}

func GetFirstDateOfMonth(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), 1, d.Hour(), 0, 0, 0, d.Location())
}

// Accepts yyyy-mm-dd
func SetDateFromString(date string) (time.Time, error) {
	dateString := fmt.Sprintf("%s 00:00:00", date)
	return time.Parse(time.DateTime, dateString)
}

// Accepts yyyy-mm-dd
func SetDateFromStringCurrencyStocks(date string) (time.Time, error) {
	dateString := fmt.Sprintf("%s 08:00:00", date)
	return time.Parse(time.DateTime, dateString)
}

// SG time
func SetDateTo0000Hours(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 8, 0, 0, 0, d.Location())
}
