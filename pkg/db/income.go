package db

import "time"

type Income struct {
	ID              uint      `gorm:"primaryKey"`
	TransactionDate time.Time `gorm:"type:timestamp;uniqueIndex:ix_date_type_incomes"`
	Type            string    `gorm:"uniqueIndex:ix_date_type_incomes"`
	Amount          float64
}
