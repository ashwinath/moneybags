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
