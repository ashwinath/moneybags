package db

import "time"

type Expense struct {
	ID              uint      `gorm:"primaryKey"`
	TransactionDate time.Time `gorm:"type:timestamp;uniqueIndex:ix_date_type_expenses"`
	Type            string    `gorm:"uniqueIndex:ix_date_type_expenses"`
	Amount          float64
}
