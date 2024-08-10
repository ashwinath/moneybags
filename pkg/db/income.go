package db

import (
	"github.com/ashwinath/moneybags/pkg/utils"
	"gorm.io/gorm"
)

const IncomeDatabaseName string = "income"

type Income struct {
	ID              uint           `gorm:"primaryKey"`
	TransactionDate utils.DateTime `gorm:"type:timestamp;uniqueIndex:ix_date_type_incomes" csv:"date"`
	Type            string         `gorm:"uniqueIndex:ix_date_type_incomes" csv:"type"`
	Amount          float64        `csv:"amount"`
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

// Clears the database
func (db *incomeDB) Clear() error {
	return db.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&Income{}).Error
}

// Bulk add data
func (db *incomeDB) BulkAdd(objs interface{}) error {
	return db.db.Create(objs).Error
}

func (db *incomeDB) Count() (int64, error) {
	var count int64
	r := db.db.Model(&Income{}).Count(&count)
	if r.Error != nil {
		return 0, r.Error
	}
	return count, nil
}
