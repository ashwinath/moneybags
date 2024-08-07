package db

import "time"

type Stock struct {
	ID        uint      `gorm:"primaryKey"`
	TradeDate time.Time `gorm:"type:timestamp;uniqueIndex:uidx_stocks"`
	Symbol    string    `gorm:"uniqueIndex:uidx_stocks"`
	Price     float64
}
