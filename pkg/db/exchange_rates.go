package db

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const ExchangeRateDatabaseName string = "exchange-rate"

type ExchangeRate struct {
	ID        uint      `gorm:"primaryKey"`
	TradeDate time.Time `gorm:"type:timestamptz;uniqueIndex:uidx_exchange_rates"`
	Symbol    string    `gorm:"uniqueIndex:uidx_exchange_rates"`
	Price     float64
}

type ExchangeRateDB interface {
	BulkAdd(objs any) error
	GetExchangeRateByDate(date time.Time, symbol string) (float64, error)
	GetEarliestDate(symbol string) (*time.Time, error)
}

type exchangeRateDB struct {
	db *gorm.DB
}

func NewExchangeRateDB(db *DB) (ExchangeRateDB, error) {
	if err := db.DB.AutoMigrate(&ExchangeRate{}); err != nil {
		return nil, err
	}

	return &exchangeRateDB{
		db: db.DB,
	}, nil
}

func (db *exchangeRateDB) BulkAdd(objs any) error {
	return db.db.Clauses(clause.OnConflict{DoNothing: true}).Create(objs).Error
}

func (db *exchangeRateDB) GetExchangeRateByDate(date time.Time, symbol string) (float64, error) {
	var er float64

	res := db.db.Model(ExchangeRate{}).
		Select("price").
		Where("trade_date = ?", date).
		Where("symbol = ?", symbol).
		First(&er)

	return er, res.Error
}

// GetEarliestDate returns the earliest stored trade date for the given symbol.
func (db *exchangeRateDB) GetEarliestDate(symbol string) (*time.Time, error) {
	var earliest time.Time
	res := db.db.Model(&ExchangeRate{}).
		Where("symbol = ?", symbol).
		Order("trade_date asc").
		Limit(1).
		Select("trade_date").
		Scan(&earliest)
	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, nil
	}

	return &earliest, nil
}
