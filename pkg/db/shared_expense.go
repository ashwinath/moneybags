package db

import "time"

type SharedExpense struct {
	ID          uint      `gorm:"primaryKey"`
	ExpenseDate time.Time `gorm:"type:timestamp;uniqueIndex:ix_shared_expense_type_expense_date"`
	Type        string    `gorm:"uniqueIndex:ix_shared_expense_type_expense_date"`
	Amount      float64
}

func (SharedExpense) TableName() string {
	return "shared_expense"
}
