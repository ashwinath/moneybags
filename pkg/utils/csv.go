package utils

import (
	"database/sql/driver"
	"fmt"
	"os"
	"time"

	"github.com/gocarina/gocsv"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DateTime struct {
	time.Time
}

// Convert the CSV string as internal date
func (date *DateTime) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse(time.DateOnly, csv)
	return err
}

func (DateTime) GormDataType() string {
	return "timestamp"
}

func (DateTime) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return "timestamp"
}

func (dt DateTime) Value() (driver.Value, error) {
	if !dt.IsZero() {
		return dt.GetTime().Format(time.DateTime), nil
	} else {
		return nil, nil
	}
}

func (dt *DateTime) GetTime() time.Time {
	return dt.Time
}

func UnmarshalCSV(filepath string, obj interface{}) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open file (%s) during csv unmarshalling: %s", filepath, err)
	}

	defer file.Close()

	if err := gocsv.UnmarshalFile(file, obj); err != nil { // Load clients from file
		return fmt.Errorf("failed to unmarshal file (%s) during csv unmarshalling: %s", filepath, err)
	}

	return nil
}
