package financials

import (
	"path"
	"testing"

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

	return framework.New(c, sugar, map[string]any{
		database.AssetDatabaseName:         assetDB,
		database.ExpenseDatabaseName:       expenseDB,
		database.IncomeDatabaseName:        incomeDB,
		database.SharedExpenseDatabaseName: sharedExpenseDB,
		database.TradeDatabaseName:         tradeDB,
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

func TestLoad(t *testing.T) {
	err := database.RunTest(func(db *database.DB) {
		fw := createFW(t, db)
		loader := NewLoader(fw)
		err := loader.Start()
		assert.Nil(t, err)

		allCounters := []database.Counter{
			fw.GetDB(database.AssetDatabaseName).(database.Counter),
			fw.GetDB(database.ExpenseDatabaseName).(database.Counter),
			fw.GetDB(database.IncomeDatabaseName).(database.Counter),
			fw.GetDB(database.SharedExpenseDatabaseName).(database.Counter),
			fw.GetDB(database.TradeDatabaseName).(database.Counter),
		}

		for _, counter := range allCounters {
			count, err := counter.Count()
			assert.Nil(t, err)
			assert.Greater(t, count, int64(0))
		}
	})
	assert.Nil(t, err)
}
