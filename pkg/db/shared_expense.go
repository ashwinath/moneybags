package db

import (
	"github.com/ashwinath/moneybags/pkg/utils"
	"gorm.io/gorm"
)

const SharedExpenseDatabaseName string = "shared-expense"

type SharedExpense struct {
	ID          uint           `gorm:"primaryKey"`
	ExpenseDate utils.DateTime `gorm:"type:timestamp;uniqueIndex:ix_shared_expense_type_expense_date" csv:"date"`
	Type        string         `gorm:"uniqueIndex:ix_shared_expense_type_expense_date" csv:"type"`
	Amount      float64        `csv:"amount"`
}

func (SharedExpense) TableName() string {
	return "shared_expense"
}

type SharedExpenseDB interface {
}

type sharedExpenseDB struct {
	db *gorm.DB
}

func NewSharedExpenseDB(db *DB) (SharedExpenseDB, error) {
	if err := db.DB.AutoMigrate(&SharedExpense{}); err != nil {
		return nil, err
	}

	return &sharedExpenseDB{
		db: db.DB,
	}, nil
}

// Clears the database
func (db *sharedExpenseDB) Clear() error {
	return db.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&SharedExpense{}).Error
}

// Bulk add data
func (db *sharedExpenseDB) BulkAdd(objs interface{}) error {
	return db.db.Create(objs).Error
}

func (db *sharedExpenseDB) Count() (int64, error) {
	var count int64
	r := db.db.Model(&SharedExpense{}).Count(&count)
	if r.Error != nil {
		return 0, r.Error
	}
	return count, nil
}
