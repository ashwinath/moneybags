package db

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const AverageExpenditureDatabaseName string = "average-expenditure"

type AverageExpenditure struct {
	ID          uint      `gorm:"primaryKey"`
	ExpenseDate time.Time `gorm:"type:timestamptz;unique"`
	Amount      float64
}

type AverageExpenditureDB interface {
	BulkInsertOnConflictOverride(objs []AverageExpenditure) error
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

func (db *averageExpenditureDB) BulkInsertOnConflictOverride(objs []AverageExpenditure) error {
	return db.db.Clauses(clause.OnConflict{DoNothing: true}).Create(objs).Error
}
