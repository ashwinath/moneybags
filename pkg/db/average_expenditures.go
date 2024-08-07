package db

import "time"

type AverageExpenditure struct {
	ID          uint      `gorm:"primaryKey"`
	ExpenseDate time.Time `gorm:"type:timestamp;unique"`
	Amount      float64
}
