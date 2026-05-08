package db

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const MortgageDatabaseName string = "mortgage"

type Mortgage struct {
	ID                 uint      `gorm:"primaryKey"`
	Date               time.Time `gorm:"type:timestamptz"`
	InterestPaid       float64
	PrincipalPaid      float64
	TotalInterestPaid  float64
	TotalPrincipalPaid float64
	GroupName          string
}

func (Mortgage) TableName() string {
	return "mortgage"
}

type MortgageDB interface {
	BulkUpdate([]Mortgage) error
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
func (db *mortgageDB) BulkAdd(objs any) error {
	return db.db.Clauses(clause.OnConflict{DoNothing: true}).Create(objs).Error
}

// Bulk update data
func (db *mortgageDB) BulkUpdate(mortgages []Mortgage) error {
	return db.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"interest_paid",
			"principal_paid",
			"total_interest_paid",
			"total_principal_paid",
			"group_name",
		}),
	}).Create(mortgages).Error
}

func (db *mortgageDB) GetMortgage() ([]Mortgage, error) {
	mortgages := []Mortgage{}
	res := db.db.Order("date asc").Find(&mortgages)
	return mortgages, res.Error
}
