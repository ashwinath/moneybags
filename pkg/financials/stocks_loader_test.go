package financials

import (
	"testing"

	database "github.com/ashwinath/moneybags/pkg/db"
	"github.com/stretchr/testify/assert"
)

func TestStocksLoader(t *testing.T) {
	err := database.RunTest(func(db *database.DB) {
		fw := createFW(t, db)

		// Load trades first
		loader := NewCSVLoader(fw)
		err := loader.Load()
		assert.Nil(t, err)

		av := NewFakeAlphavantage()
		stocksLoader := NewStocksLoader(fw, av)
		err = stocksLoader.Load()
		assert.Nil(t, err)

		var stockCount int64
		res := db.DB.Model(database.Symbol{}).Where(
			"symbol_type = ?", database.SymbolTypeStock,
		).Count(&stockCount)
		assert.Nil(t, res.Error)
		assert.Equal(t, int64(3), stockCount)

		var currencyCount int64
		res = db.DB.Model(database.Symbol{}).Where(
			"symbol_type = ?", database.SymbolTypeCurrency,
		).Count(&currencyCount)
		assert.Nil(t, res.Error)
		assert.Equal(t, int64(1), currencyCount)

		var exchangeRateCount int64
		res = db.DB.Model(database.ExchangeRate{}).Count(&exchangeRateCount)
		assert.Nil(t, res.Error)
		assert.Equal(t, int64(1), currencyCount)

		currencySymbol := database.Symbol{}
		res = db.DB.Model(database.Symbol{}).
			Where("symbol_type = ?", database.SymbolTypeCurrency).
			First(&currencySymbol)
		assert.Nil(t, res.Error)
		assert.NotNil(t, currencySymbol.LastProcessedDate)

		stockSymbol := database.Symbol{}
		res = db.DB.Model(database.Symbol{}).
			Where("symbol_type = ?", database.SymbolTypeStock).
			First(&stockSymbol)
		assert.Nil(t, res.Error)
		assert.NotNil(t, stockSymbol.LastProcessedDate)
	})

	assert.Nil(t, err)
}
