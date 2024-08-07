package db

import "time"

type ExchangeRate struct {
	ID        uint      `gorm:"primaryKey"`
	TradeDate time.Time `gorm:"type:timestamp;uniqueIndex:uidx_exchange_rates"`
	Symbol    string    `gorm:"uniqueIndex:uidx_exchange_rates"`
	Amount    float64
}
