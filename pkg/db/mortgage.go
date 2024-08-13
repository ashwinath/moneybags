package db

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const MortgageDatabaseName string = "mortgage"

type Mortgage struct {
	ID                 uint      `gorm:"primaryKey"`
	Date               time.Time `gorm:"type:timestamp;unique"`
	InterestPaid       float64
	PrincipalPaid      float64
	TotalInterestPaid  float64
	TotalPrincipalPaid float64
	TotalInterestLeft  float64
	TotalPrincipalLeft float64
}

func (Mortgage) TableName() string {
	return "mortgage"
}

type MortgageDB interface {
	GetMortgage() ([]Mortgage, error)
}

type mortgageDB struct {
	db *gorm.DB
}

func NewMortgageDB(db *DB) (MortgageDB, error) {
	if err := db.DB.AutoMigrate(&Mortgage{}); err != nil {
		return nil, err
	}

	return &mortgageDB{
		db: db.DB,
	}, nil
}

// Clears the database
func (db *mortgageDB) Clear() error {
	return db.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&Mortgage{}).Error
}

// Bulk add data
func (db *mortgageDB) BulkAdd(objs interface{}) error {
	return db.db.Clauses(clause.OnConflict{DoNothing: true}).Create(objs).Error
}

func (db *mortgageDB) GetMortgage() ([]Mortgage, error) {
	mortgages := []Mortgage{}
	res := db.db.Order("date asc").Find(&mortgages)
	return mortgages, res.Error
}
