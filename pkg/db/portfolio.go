package db

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const PortfolioDatabaseName string = "portfolio"

type Portfolio struct {
	ID            uint      `gorm:"primaryKey"`
	TradeDate     time.Time `gorm:"type:timestamp;uniqueIndex:uidx_portfolios"`
	Symbol        string    `gorm:"uniqueIndex:uidx_portfolios"`
	Principal     float64
	NAV           float64
	SimpleReturns float64
	Quantity      float64
}

type PortfolioDB interface {
	BulkAdd(objs interface{}) error
	GetFirstTradeDate() (time.Time, error)
	GetPortfolioAmountByDate(date time.Time) (float64, error)
}

type portfolioDB struct {
	db *gorm.DB
}

func NewPortfolioDB(db *DB) (PortfolioDB, error) {
	if err := db.DB.AutoMigrate(&Portfolio{}); err != nil {
		return nil, err
	}

	return &portfolioDB{
		db: db.DB,
	}, nil
}

// Bulk add data
func (db *portfolioDB) BulkAdd(objs interface{}) error {
	return db.db.Clauses(clause.OnConflict{DoNothing: true}).Create(objs).Error
}

func (db *portfolioDB) GetFirstTradeDate() (time.Time, error) {
	var val time.Time
	res := db.db.Model(Portfolio{}).
		Select("trade_date").
		Order("trade_date asc").
		First(&val)
	return val, res.Error
}

func (db *portfolioDB) GetPortfolioAmountByDate(date time.Time) (float64, error) {
	var val float64
	res := db.db.Model(Portfolio{}).
		Select("sum(nav)").
		Where("trade_date = ?", date).
		Scan(&val)
	return val, res.Error
}
