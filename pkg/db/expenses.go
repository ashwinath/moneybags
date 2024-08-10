package db

import (
	"github.com/ashwinath/moneybags/pkg/utils"
	"gorm.io/gorm"
)

const ExpenseDatabaseName string = "expense"

type Expense struct {
	ID              uint           `gorm:"primaryKey"`
	TransactionDate utils.DateTime `gorm:"type:timestamp;uniqueIndex:ix_date_type_expenses" csv:"date"`
	Type            string         `gorm:"uniqueIndex:ix_date_type_expenses" csv:"type"`
	Amount          float64        `csv:"amount"`
}

type ExpenseDB interface {
	BulkAdd(objs interface{}) error
}

type expenseDB struct {
	db *gorm.DB
}

func NewExpenseDB(db *DB) (ExpenseDB, error) {
	if err := db.DB.AutoMigrate(&Expense{}); err != nil {
		return nil, err
	}

	return &expenseDB{
		db: db.DB,
	}, nil
}

// Clears the database
func (db *expenseDB) Clear() error {
	return db.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&Expense{}).Error
}

// Bulk add data
func (db *expenseDB) BulkAdd(objs interface{}) error {
	return db.db.Create(objs).Error
}

func (db *expenseDB) Count() (int64, error) {
	var count int64
	r := db.db.Model(&Expense{}).Count(&count)
	if r.Error != nil {
		return 0, r.Error
	}
	return count, nil
}
