package financials

import (
	"path"
	"testing"
	"time"

	"github.com/ashwinath/moneybags/pbgo/configpb"
	"github.com/ashwinath/moneybags/pkg/config"
	database "github.com/ashwinath/moneybags/pkg/db"
	"github.com/ashwinath/moneybags/pkg/framework"
	"github.com/ashwinath/moneybags/pkg/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func createFW(t *testing.T, db *database.DB) framework.FW {
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()

	p := path.Join(utils.GetLocalRepoLocation(), "./pkg/config/testdata/config.yaml")
	c, err := config.New(p)
	assert.Nil(t, err)

	subsituteLocalRepoLocation(c)

	assetDB, err := database.NewAssetDB(db)
	assert.Nil(t, err)

	expenseDB, err := database.NewExpenseDB(db)
	assert.Nil(t, err)

	incomeDB, err := database.NewIncomeDB(db)
	assert.Nil(t, err)

	sharedExpenseDB, err := database.NewSharedExpenseDB(db)
	assert.Nil(t, err)

	tradeDB, err := database.NewTradeDB(db)
	assert.Nil(t, err)

	txDB, err := database.NewTransactionDB(db)
	assert.Nil(t, err)

	symbolsDB, err := database.NewSymbolDB(db)
	assert.Nil(t, err)

	exchangeRateDB, err := database.NewExchangeRateDB(db)
	assert.Nil(t, err)

	stockDB, err := database.NewStockDB(db)
	assert.Nil(t, err)

	return framework.New(c, sugar, map[string]any{
		database.AssetDatabaseName:         assetDB,
		database.ExpenseDatabaseName:       expenseDB,
		database.IncomeDatabaseName:        incomeDB,
		database.SharedExpenseDatabaseName: sharedExpenseDB,
		database.TradeDatabaseName:         tradeDB,
		database.TransactionDatabaseName:   txDB,
		database.SymbolDatabaseName:        symbolsDB,
		database.ExchangeRateDatabaseName:  exchangeRateDB,
		database.StockDatabaseName:         stockDB,
	})
}

func subsituteLocalRepoLocation(c *configpb.Config) {
	p := c.FinancialsData.AssetsCsvFilepath
	newPath := path.Join(utils.GetLocalRepoLocation(), p)
	c.FinancialsData.AssetsCsvFilepath = newPath

	p = c.FinancialsData.ExpensesCsvFilepath
	newPath = path.Join(utils.GetLocalRepoLocation(), p)
	c.FinancialsData.ExpensesCsvFilepath = newPath

	p = c.FinancialsData.IncomeCsvFilepath
	newPath = path.Join(utils.GetLocalRepoLocation(), p)
	c.FinancialsData.IncomeCsvFilepath = newPath

	p = c.FinancialsData.SharedExpensesCsvFilepath
	newPath = path.Join(utils.GetLocalRepoLocation(), p)
	c.FinancialsData.SharedExpensesCsvFilepath = newPath

	p = c.FinancialsData.TradesCsvFilepath
	newPath = path.Join(utils.GetLocalRepoLocation(), p)
	c.FinancialsData.TradesCsvFilepath = newPath
}

func parseDateForced(t *testing.T, dateString string) time.Time {
	loc, err := time.LoadLocation("Asia/Singapore")
	assert.Nil(t, err)

	parsed, err := time.ParseInLocation(time.DateOnly, dateString, loc)
	assert.Nil(t, err)

	return parsed
}

type fakeAlphavantage struct{}

func (fakeAlphavantage) GetSymbolFromAlphavantage(symbol string) (*AlphavantageSymbol, error) {
	return &AlphavantageSymbol{
		Symbol:   symbol,
		Currency: "USD",
	}, nil
}

func (fakeAlphavantage) GetCurrencyHistory(from string, to string, isCompact bool) (map[string]OHLC, error) {
	return map[string]OHLC{
		"2021-08-19": {
			Open:  1.00,
			High:  1.00,
			Low:   1.00,
			Close: 1.00,
		},
		"2021-08-30": {
			Open:  1.00,
			High:  1.00,
			Low:   1.00,
			Close: 1.00,
		},
	}, nil
}

func (fakeAlphavantage) GetStockHistory(symbol string, isCompact bool) (map[string]OHLC, error) {
	return map[string]OHLC{
		"2021-08-19": {
			Open:  2.00,
			High:  2.00,
			Low:   2.00,
			Close: 2.00,
		},
		"2021-08-30": {
			Open:  2.00,
			High:  2.00,
			Low:   2.00,
			Close: 2.00,
		},
	}, nil
}

func NewFakeAlphavantage() Alphavantage {
	return fakeAlphavantage{}
}
