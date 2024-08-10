package db

import (
	"testing"
	"time"

	"github.com/ashwinath/moneybags/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetUniqueSymbols(t *testing.T) {
	err := RunTest(func(db *DB) {
		tradeDB, err := NewTradeDB(db)
		assert.Nil(t, err)

		trades := []*Trade{
			{
				DatePurchased: utils.DateTime{Time: time.Now()},
				Symbol:        "hello",
				PriceEach:     100.23,
				Quantity:      42,
				TradeType:     "buy",
			},
			{
				DatePurchased: utils.DateTime{Time: time.Now()},
				Symbol:        "world",
				PriceEach:     100.23,
				Quantity:      42,
				TradeType:     "buy",
			},
			{
				DatePurchased: utils.DateTime{Time: time.Now()},
				Symbol:        "world",
				PriceEach:     100.23,
				Quantity:      42,
				TradeType:     "buy",
			},
		}
		tradeBulkerDB := tradeDB.(ClearAndBulkAdder)
		err = tradeBulkerDB.BulkAdd(&trades)
		assert.Nil(t, err)

		// query
		symbols, err := tradeDB.GetUniqueSymbols()
		assert.Nil(t, err)

		assert.Len(t, symbols, 2)
	})
	assert.Nil(t, err)
}
