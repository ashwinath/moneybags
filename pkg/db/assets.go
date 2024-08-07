package db

import "time"

type Assets struct {
	ID              uint      `gorm:"primaryKey"`
	TransactionDate time.Time `gorm:"type:timestamp;uniqueIndex:ix_date_type_assets"`
	Type            string    `gorm:"uniqueIndex:ix_date_type_assets"`
	Amount          float64
}
