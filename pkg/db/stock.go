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
	BulkAdd(objs interface{}) error
	GetStockPrice(date time.Time, symbol string) (float64, error)
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

func (db *stockDB) BulkAdd(objs interface{}) error {
	return db.db.Create(objs).Error
}

func (db *stockDB) GetStockPrice(date time.Time, symbol string) (float64, error) {
	var val float64

	res := db.db.Model(Stock{}).
		Select("price").
		Where("trade_date = ?", date).
		Where("symbol = ?", symbol).
		First(&val)

	return val, res.Error
}
