package db

import (
	"github.com/ashwinath/moneybags/pkg/utils"
	"gorm.io/gorm"
)

const TradeDatabaseName string = "trade"

type Trade struct {
	ID            uint           `gorm:"primaryKey"`
	DatePurchased utils.DateTime `gorm:"type:timestamp" csv:"date_purchased"`
	Symbol        string         `csv:"symbol"`
	PriceEach     float64        `csv:"price_each"`
	Quantity      float64        `csv:"quantity"`
	TradeType     string         `csv:"trade_type"`
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

// Clears the database
func (db *tradeDB) Clear() error {
	return db.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&Trade{}).Error
}

// Bulk add data
func (db *tradeDB) BulkAdd(objs interface{}) error {
	return db.db.Create(objs).Error
}

func (db *tradeDB) Count() (int64, error) {
	var count int64
	r := db.db.Model(&Trade{}).Count(&count)
	if r.Error != nil {
		return 0, r.Error
	}
	return count, nil
}
