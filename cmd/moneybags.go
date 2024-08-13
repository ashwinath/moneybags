package main

import (
	"context"
	"flag"

	"github.com/ashwinath/moneybags/pkg/config"
	"github.com/ashwinath/moneybags/pkg/db"
	"github.com/ashwinath/moneybags/pkg/financials"
	"github.com/ashwinath/moneybags/pkg/framework"
	"github.com/ashwinath/moneybags/pkg/modules"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	sugar.Info("Starting moneybags")

	sugar.Info("Loading config")

	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	c, err := config.New(*configPath)
	if err != nil {
		sugar.Fatalf("Failed to load config, %v", err)
	}

	// load databases
	baseDB, err := db.NewBaseDB(c.PostgresDb)
	if err != nil {
		sugar.Fatalf("Failed to initialise Base DB, %v", err)
	}
	defer baseDB.Close()

	// Load framework
	fw := framework.New(c, sugar, createDBs(baseDB, sugar))

	// Load modules
	telegram, err := modules.NewTelegramModule(fw)
	if err != nil {
		sugar.Fatalf("Failed to initialise telegram module, %v", err)
	}

	financials, err := modules.NewFinancialsModule(fw, financials.NewAlphavantage(
		c.FinancialsConfig.AlphavantageApiKey,
	))
	if err != nil {
		sugar.Fatalf("Failed to initialise financials module, %v", err)
	}

	// Run app
	app := framework.NewApp(sugar, telegram, financials)

	ctx, cancel := context.WithCancel(context.Background())
	framework.ListenForSignal(cancel, sugar)

	app.Run(ctx)
}

func createDBs(baseDB *db.DB, sugar *zap.SugaredLogger) map[string]any {
	assetDB, err := db.NewAssetDB(baseDB)
	if err != nil {
		sugar.Fatalf("Failed to initialise Asset DB, %v", err)
	}

	averageExpenditureDB, err := db.NewAverageExpenditureDB(baseDB)
	if err != nil {
		sugar.Fatalf("Failed to initialise Average Expenditure DB, %v", err)
	}

	exchangeRateDB, err := db.NewExchangeRateDB(baseDB)
	if err != nil {
		sugar.Fatalf("Failed to initialise Exchange Rate DB, %v", err)
	}

	expenseDB, err := db.NewExpenseDB(baseDB)
	if err != nil {
		sugar.Fatalf("Failed to initialise Expense DB, %v", err)
	}

	incomeDB, err := db.NewIncomeDB(baseDB)
	if err != nil {
		sugar.Fatalf("Failed to initialise Income DB, %v", err)
	}

	mortgageDB, err := db.NewMortgageDB(baseDB)
	if err != nil {
		sugar.Fatalf("Failed to initialise Mortgage DB, %v", err)
	}

	portfolioDB, err := db.NewPortfolioDB(baseDB)
	if err != nil {
		sugar.Fatalf("Failed to initialise Portfolio DB, %v", err)
	}

	sharedExpenseDB, err := db.NewSharedExpenseDB(baseDB)
	if err != nil {
		sugar.Fatalf("Failed to initialise Expense DB, %v", err)
	}

	stockDB, err := db.NewStockDB(baseDB)
	if err != nil {
		sugar.Fatalf("Failed to initialise Stock DB, %v", err)
	}

	symbolDB, err := db.NewSymbolDB(baseDB)
	if err != nil {
		sugar.Fatalf("Failed to initialise Symbol DB, %v", err)
	}

	tradeDB, err := db.NewTradeDB(baseDB)
	if err != nil {
		sugar.Fatalf("Failed to initialise Trade DB, %v", err)
	}

	transactionDB, err := db.NewTransactionDB(baseDB)
	if err != nil {
		sugar.Fatalf("Failed to initialise Transaction DB, %v", err)
	}

	return map[string]any{
		db.AssetDatabaseName:              assetDB,
		db.AverageExpenditureDatabaseName: averageExpenditureDB,
		db.ExchangeRateDatabaseName:       exchangeRateDB,
		db.ExpenseDatabaseName:            expenseDB,
		db.IncomeDatabaseName:             incomeDB,
		db.MortgageDatabaseName:           mortgageDB,
		db.PortfolioDatabaseName:          portfolioDB,
		db.SharedExpenseDatabaseName:      sharedExpenseDB,
		db.StockDatabaseName:              stockDB,
		db.SymbolDatabaseName:             symbolDB,
		db.TradeDatabaseName:              tradeDB,
		db.TransactionDatabaseName:        transactionDB,
	}
}
