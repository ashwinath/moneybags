package db

import "time"

type Trade struct {
	ID            uint      `gorm:"primaryKey"`
	DatePurchased time.Time `gorm:"type:timestamp"`
	Symbol        string
	PriceEach     float64
	Quantity      float64
	TradeType     string
}
