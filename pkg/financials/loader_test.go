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

		// stock loader
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

		var portfolioCount int64
		res = db.DB.Model(database.Portfolio{}).Count(&portfolioCount)
		assert.Nil(t, res.Error)
		assert.Greater(t, portfolioCount, int64(1))

		// investment loader
		investmentsLoader := NewInvestmentsLoader(fw)
		err = investmentsLoader.Load()
		assert.Nil(t, err)

		assets := []database.Asset{}
		res = db.DB.Where("type = ?", "Investments").Find(&assets)
		assert.Nil(t, res.Error)
		assert.Greater(t, len(assets), 1)
	})

	assert.Nil(t, err)
}
