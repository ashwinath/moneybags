package db

import (
	"time"

	"gorm.io/gorm"
)

const StockDatabaseName string = "stock"

type Stock struct {
	ID        uint      `gorm:"primaryKey"`
	TradeDate time.Time `gorm:"type:timestamp;uniqueIndex:uidx_stocks"`
	Symbol    string    `gorm:"uniqueIndex:uidx_stocks"`
	Price     float64
}

type StockDB interface {
}

type stockDB struct {
	db *gorm.DB
}

func NewStockDB(db *DB) (StockDB, error) {
	if err := db.DB.AutoMigrate(&Stock{}); err != nil {
		return nil, err
	}

	return &stockDB{
		db: db.DB,
	}, nil
}
