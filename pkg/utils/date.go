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
