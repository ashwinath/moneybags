package db

import (
	"time"

	"gorm.io/gorm"
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
