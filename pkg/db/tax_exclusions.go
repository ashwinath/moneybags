package db

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const TaxExclusions string = "tax_exclusion"

type TaxExclusion struct {
	ID     uint    `gorm:"primaryKey"`
	Year   uint    `csv:"year"`
	Type   string  `csv:"type"`
	Amount float64 `csv:"amount"`
}

type taxExlusionDB struct {
	db *gorm.DB
}

func NewTaxExclusionDB(db *DB) (ClearAndBulkAdder, error) {
	if err := db.DB.AutoMigrate(&TaxExclusion{}); err != nil {
		return nil, err
	}

	return &taxExlusionDB{
		db: db.DB,
	}, nil
}

// Clears the database
func (db *taxExlusionDB) Clear() error {
	return db.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&Asset{}).Error
}

// Bulk add data
func (db *taxExlusionDB) BulkAdd(objs any) error {
	return db.db.Clauses(clause.OnConflict{DoNothing: true}).Create(objs).Error
}

// Counter interface
func (db *taxExlusionDB) Count() (int64, error) {
	var count int64
	r := db.db.Model(&TaxExclusion{}).Count(&count)
	if r.Error != nil {
		return 0, r.Error
	}
	return count, nil
}
