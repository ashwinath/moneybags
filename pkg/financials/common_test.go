package financials

import (
	"path"
	"testing"
	"time"

	"github.com/ashwinath/moneybags/pbgo/configpb"
	"github.com/ashwinath/moneybags/pkg/config"
	"github.com/ashwinath/moneybags/pkg/db"
	database "github.com/ashwinath/moneybags/pkg/db"
	"github.com/ashwinath/moneybags/pkg/framework"
	"github.com/ashwinath/moneybags/pkg/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func createFW(t *testing.T, baseDB *database.DB) framework.FW {
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()

	p := path.Join(utils.GetLocalRepoLocation(), "./pkg/config/testdata/config.yaml")
	c, err := config.New(p)
	assert.Nil(t, err)

	subsituteLocalRepoLocation(c)

	assetDB, err := db.NewAssetDB(baseDB)
	assert.Nil(t, err)

	averageExpenditureDB, err := db.NewAverageExpenditureDB(baseDB)
	assert.Nil(t, err)

	exchangeRateDB, err := db.NewExchangeRateDB(baseDB)
	assert.Nil(t, err)

	expenseDB, err := db.NewExpenseDB(baseDB)
	assert.Nil(t, err)

	incomeDB, err := db.NewIncomeDB(baseDB)
	assert.Nil(t, err)

	mortgageDB, err := db.NewMortgageDB(baseDB)
	assert.Nil(t, err)

	portfolioDB, err := db.NewPortfolioDB(baseDB)
	assert.Nil(t, err)

	sharedExpenseDB, err := db.NewSharedExpenseDB(baseDB)
	assert.Nil(t, err)

	stockDB, err := db.NewStockDB(baseDB)
	assert.Nil(t, err)

	symbolDB, err := db.NewSymbolDB(baseDB)
	assert.Nil(t, err)

	tradeDB, err := db.NewTradeDB(baseDB)
	assert.Nil(t, err)

	transactionDB, err := db.NewTransactionDB(baseDB)
	assert.Nil(t, err)

	return framework.New(c, sugar, map[string]any{
		database.AssetDatabaseName:              assetDB,
		database.AverageExpenditureDatabaseName: averageExpenditureDB,
		database.ExchangeRateDatabaseName:       exchangeRateDB,
		database.ExpenseDatabaseName:            expenseDB,
		database.IncomeDatabaseName:             incomeDB,
		database.MortgageDatabaseName:           mortgageDB,
		database.PortfolioDatabaseName:          portfolioDB,
		database.SharedExpenseDatabaseName:      sharedExpenseDB,
		database.StockDatabaseName:              stockDB,
		database.SymbolDatabaseName:             symbolDB,
		database.TradeDatabaseName:              tradeDB,
		database.TransactionDatabaseName:        transactionDB,
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
	return autoGenOHLC(), nil
}

func (fakeAlphavantage) GetStockHistory(symbol string, isCompact bool) (map[string]OHLC, error) {
	return autoGenOHLC(), nil
}

func autoGenOHLC() map[string]OHLC {
	ret := map[string]OHLC{}
	date, _ := time.Parse(time.DateOnly, "2021-08-19")
	for date.Before(time.Now().AddDate(0, 0, 1)) {
		ret[date.Format(time.DateOnly)] = OHLC{Open: 2.00, High: 2.00, Low: 2.00, Close: 2.00}
		date = date.AddDate(0, 0, 1)
	}
	return ret
}

func NewFakeAlphavantage() Alphavantage {
	return fakeAlphavantage{}
}
