package financials

import (
	"testing"
	"time"

	database "github.com/ashwinath/moneybags/pkg/db"
	"github.com/ashwinath/moneybags/pkg/utils"
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
		provider := NewFakeMarketDataProvider()
		stocksLoader := NewStocksLoader(fw, provider)
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
		res = db.DB.Where("type = ?", "Investments").Order("transaction_date asc").Find(&assets)
		assert.Nil(t, res.Error)
		assert.Greater(t, len(assets), 1)

		// verify the Investments assets exclude CSPX.L (configured in test config)
		// by mirroring the loader's date computation and comparing amounts
		portfolioDB := fw.GetDB(database.PortfolioDatabaseName).(database.PortfolioDB)
		firstDate, err := portfolioDB.GetFirstTradeDate()
		assert.Nil(t, err)

		currentDate := firstDate
		if currentDate.Day() != 1 {
			currentDate = time.Date(currentDate.Year(), currentDate.Month(), 1, currentDate.Hour(), 0, 0, 0, currentDate.Location())
			currentDate = currentDate.AddDate(0, 1, 0)
			currentDate = utils.SetDateTo0000Hours(currentDate)
		}

		for i, asset := range assets[:3] {
			fullAmount, err := portfolioDB.GetPortfolioAmountByDate(currentDate, nil)
			assert.Nil(t, err)

			excludedAmount, err := portfolioDB.GetPortfolioAmountByDate(currentDate, []string{"CSPX.L"})
			assert.Nil(t, err)

			assert.Greater(t, fullAmount, excludedAmount)
			assert.Equal(t, excludedAmount, asset.Amount, "asset row %d mismatch", i)

			currentDate = currentDate.AddDate(0, 1, 0)
		}

		// shared expense loader
		sharedExpenseLoader := NewSharedExpenseLoader(fw)
		err = sharedExpenseLoader.Load()
		assert.Nil(t, err)

		nonSpecialExpenses := []database.Expense{}
		res = db.DB.Where("type = ?", nonSpecialSharedExpensesType).Find(&nonSpecialExpenses)
		assert.Nil(t, res.Error)
		assert.Greater(t, len(nonSpecialExpenses), 1)

		specialExpenses := []database.Expense{}
		res = db.DB.Where("type = ?", specialSharedExpensesType).Find(&specialExpenses)
		assert.Nil(t, res.Error)
		assert.Greater(t, len(specialExpenses), 1)

		// average expenditure loader
		averageExpenditureLoader := NewAverageExpenditureLoader(fw)
		err = averageExpenditureLoader.Load()
		assert.Nil(t, err)

		averageExpenditure := []database.AverageExpenditure{}
		res = db.DB.Find(&averageExpenditure)
		assert.Nil(t, res.Error)
		assert.Greater(t, len(averageExpenditure), 1)
	})

	assert.Nil(t, err)
}

func TestStocksLoaderCleansUpLegacySymbols(t *testing.T) {
	err := database.RunTest(func(db *database.DB) {
		fw := createFW(t, db)

		legacySymbol := "AAAA.LON"
		startDate, err := time.Parse(time.DateOnly, "2021-08-19")
		assert.Nil(t, err)

		items := []database.Symbol{
			{SymbolType: database.SymbolTypeStock, Symbol: legacySymbol, BaseCurrency: &[]string{"USD"}[0]},
			{SymbolType: database.SymbolTypeCurrency, Symbol: "USD"},
		}
		assert.Nil(t, db.DB.Create(&items).Error)

		assert.Nil(t, db.DB.Create([]database.Stock{
			{TradeDate: startDate, Symbol: legacySymbol, Price: 2.0},
		}).Error)

		assert.Nil(t, db.DB.Create([]database.Portfolio{
			{TradeDate: startDate, Symbol: legacySymbol, Principal: 20.0, NAV: 20.0, SimpleReturns: 0.0, Quantity: 10.0},
		}).Error)

		loader := NewCSVLoader(fw)
		assert.Nil(t, loader.Load())

		stocksLoader := NewStocksLoader(fw, NewFakeMarketDataProvider())
		assert.Nil(t, stocksLoader.Load())

		var count int64
		res := db.DB.Model(database.Symbol{}).Where("symbol = ?", legacySymbol).Count(&count)
		assert.Nil(t, res.Error)
		assert.Equal(t, int64(0), count)

		res = db.DB.Model(database.Stock{}).Where("symbol = ?", legacySymbol).Count(&count)
		assert.Nil(t, res.Error)
		assert.Equal(t, int64(0), count)

		res = db.DB.Model(database.Portfolio{}).Where("symbol = ?", legacySymbol).Count(&count)
		assert.Nil(t, res.Error)
		assert.Equal(t, int64(0), count)
	})

	assert.Nil(t, err)
}
