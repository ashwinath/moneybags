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
