package db

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const TotalCompensationName string = "total_compensation"

type TotalCompensation struct {
	Year              int `gorm:"primaryKey"`
	GrossBase         float64
	GrossBonus        float64
	RSU               float64
	ESPP              float64
	GrossTC           float64
	CPF               float64
	TotalCompensation float64
}

type TotalCompensationDB any

type totalCompensationDB struct {
	db *gorm.DB
}

func NewTotalCompensationDB(db *DB) (TotalCompensationDB, error) {
	if err := db.DB.AutoMigrate(&TotalCompensation{}); err != nil {
		return nil, err
	}

	return &totalCompensationDB{
		db: db.DB,
	}, nil
}

// Clears the database
func (db *totalCompensationDB) Clear() error {
	return db.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&TotalCompensation{}).Error
}

// Bulk add data
func (db *totalCompensationDB) BulkAdd(objs any) error {
	return db.db.Clauses(clause.OnConflict{DoNothing: true}).Create(objs).Error
}
