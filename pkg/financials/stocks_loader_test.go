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
	})

	assert.Nil(t, err)
}
