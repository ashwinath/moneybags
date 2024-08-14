package db

import (
	"time"

	"github.com/ashwinath/moneybags/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const ExpenseDatabaseName string = "expense"

type Expense struct {
	ID              uint           `gorm:"primaryKey"`
	TransactionDate utils.DateTime `gorm:"type:timestamptz;uniqueIndex:ix_date_type_expenses" csv:"date"`
	Type            string         `gorm:"uniqueIndex:ix_date_type_expenses" csv:"type"`
	Amount          float64        `csv:"amount"`
}

type ExpenseDB interface {
	BulkAdd(objs interface{}) error
	GetFirstDate() (time.Time, error)
	GetYearlyExpense(exclusionTypes []string, currentDate time.Time, windowPeriod int) (float64, error)
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
	return db.db.Clauses(clause.OnConflict{DoNothing: true}).Create(objs).Error
}

func (db *expenseDB) Count() (int64, error) {
	var count int64
	r := db.db.Model(&Expense{}).Count(&count)
	if r.Error != nil {
		return 0, r.Error
	}
	return count, nil
}

func (db *expenseDB) GetFirstDate() (time.Time, error) {
	var val time.Time
	res := db.db.Model(Expense{}).
		Select("transaction_date").
		Order("transaction_date asc").
		First(&val)
	return val, res.Error
}

func (db *expenseDB) GetYearlyExpense(exclusionTypes []string, currentDate time.Time, windowPeriod int) (float64, error) {
	query := db.db.Model(Expense{}).
		Select("sum(amount)").
		Where("transaction_date > ?", currentDate.AddDate(0, -windowPeriod, 0)).
		Where("transaction_date <= ?", currentDate)

	for _, et := range exclusionTypes {
		query = query.Where("type not like ?", et)
	}

	var val float64
	res := query.Scan(&val)
	return val, res.Error
}
