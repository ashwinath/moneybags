package db

import (
	"github.com/ashwinath/moneybags/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const IncomeDatabaseName string = "income"

type Income struct {
	ID              uint           `gorm:"primaryKey"`
	TransactionDate utils.DateTime `gorm:"type:timestamptz;uniqueIndex:ix_date_type_incomes" csv:"date"`
	Type            string         `gorm:"uniqueIndex:ix_date_type_incomes" csv:"type"`
	Amount          float64        `csv:"amount"`
}

type IncomeDB interface {
	AggregateByYear([]string) ([]AggregateIncomeByYear, error)
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
func (db *incomeDB) BulkAdd(objs any) error {
	return db.db.Clauses(clause.OnConflict{DoNothing: true}).Create(objs).Error
}

func (db *incomeDB) Count() (int64, error) {
	var count int64
	r := db.db.Model(&Income{}).Count(&count)
	if r.Error != nil {
		return 0, r.Error
	}
	return count, nil
}

type AggregateIncomeByYear struct {
	Year   int
	Amount float64
}

func (db *incomeDB) AggregateByYear(types []string) ([]AggregateIncomeByYear, error) {
	results := []AggregateIncomeByYear{}
	err := db.db.Table("incomes").
		Select("date_part('year', date_trunc('year', transaction_date)) as year, sum(amount) as amount").
		Where("type in ?", types).
		Group("year").
		Scan(&results).
		Error

	if err != nil {
		return nil, err
	}
	return results, nil
}
