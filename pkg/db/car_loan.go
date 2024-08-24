package db

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const CarLoansDatabaseName string = "car_loan"

type CarLoan struct {
	ID         uint      `gorm:"primaryKey"`
	Date       time.Time `gorm:"type:timestamptz;uniqueIndex:uidx_car_loan"`
	Name       string    `gorm:"uniqueIndex:uidx_car_loan"`
	AmountPaid float64
	AmountLeft float64
}

type CarLoanDB interface {
	Clear() error
	BulkAdd(objs interface{}) error
}

type carLoanDB struct {
	db *gorm.DB
}

func NewCarLoanDB(db *DB) (CarLoanDB, error) {
	if err := db.DB.AutoMigrate(&CarLoan{}); err != nil {
		return nil, err
	}

	return &carLoanDB{
		db: db.DB,
	}, nil
}

// Clears the database
func (db *carLoanDB) Clear() error {
	return db.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&CarLoan{}).Error
}

// Bulk add data
func (db *carLoanDB) BulkAdd(objs interface{}) error {
	return db.db.Clauses(clause.OnConflict{DoNothing: true}).Create(objs).Error
}
