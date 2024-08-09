package db

import (
	"time"

	"gorm.io/gorm"
)

const ExchangeRateDatabaseName string = "exchange-rate"

type ExchangeRate struct {
	ID        uint      `gorm:"primaryKey"`
	TradeDate time.Time `gorm:"type:timestamp;uniqueIndex:uidx_exchange_rates"`
	Symbol    string    `gorm:"uniqueIndex:uidx_exchange_rates"`
	Amount    float64
}

type ExchangeRateDB interface {
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
