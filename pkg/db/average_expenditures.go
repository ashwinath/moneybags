package db

import (
	"time"

	"gorm.io/gorm"
)

const AverageExpenditureDatabaseName string = "average-expenditure"

type AverageExpenditure struct {
	ID          uint      `gorm:"primaryKey"`
	ExpenseDate time.Time `gorm:"type:timestamp;unique"`
	Amount      float64
}

type AverageExpenditureDB interface {
}

type averageExpenditureDB struct {
	db *gorm.DB
}

func NewAverageExpenditureDB(db *DB) (AverageExpenditureDB, error) {
	if err := db.DB.AutoMigrate(&AverageExpenditure{}); err != nil {
		return nil, err
	}

	return &averageExpenditureDB{
		db: db.DB,
	}, nil
}
