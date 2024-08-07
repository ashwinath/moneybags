package db

import "time"

type Portfolio struct {
	ID            uint      `gorm:"primaryKey"`
	TradeDate     time.Time `gorm:"type:timestamp;uniqueIndex:uidx_portfolios"`
	Symbol        string    `gorm:"uniqueIndex:uidx_portfolios"`
	Principal     float64
	NAV           float64
	SimpleReturns float64
	Quantity      float64
}
