package db

import (
	"time"

	"gorm.io/gorm"
)

const ExpenseDatabaseName string = "expense"

type Expense struct {
	ID              uint      `gorm:"primaryKey"`
	TransactionDate time.Time `gorm:"type:timestamp;uniqueIndex:ix_date_type_expenses"`
	Type            string    `gorm:"uniqueIndex:ix_date_type_expenses"`
	Amount          float64
}

type ExpenseDB interface {
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
