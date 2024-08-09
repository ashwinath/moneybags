package db

import (
	"time"

	"gorm.io/gorm"
)

const TradeDatabaseName string = "trade"

type Trade struct {
	ID            uint      `gorm:"primaryKey"`
	DatePurchased time.Time `gorm:"type:timestamp"`
	Symbol        string
	PriceEach     float64
	Quantity      float64
	TradeType     string
}

type TradeDB interface {
}

type tradeDB struct {
	db *gorm.DB
}

func NewTradeDB(db *DB) (TradeDB, error) {
	if err := db.DB.AutoMigrate(&Trade{}); err != nil {
		return nil, err
	}

	return &tradeDB{
		db: db.DB,
	}, nil
}
