package db

import (
	"time"

	"gorm.io/gorm"
)

const IncomeDatabaseName string = "income"

type Income struct {
	ID              uint      `gorm:"primaryKey"`
	TransactionDate time.Time `gorm:"type:timestamp;uniqueIndex:ix_date_type_incomes"`
	Type            string    `gorm:"uniqueIndex:ix_date_type_incomes"`
	Amount          float64
}

type IncomeDB interface {
}

type incomeDB struct {
	db *gorm.DB
}

func NewIncomeDB(db *DB) (IncomeDB, error) {
	if err := db.DB.AutoMigrate(&Income{}); err != nil {
		return nil, err
	}

	return &incomeDB{
		db: db.DB,
	}, nil
}
