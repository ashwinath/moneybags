package main

import (
	"context"
	"flag"

	"github.com/ashwinath/moneybags/pkg/config"
	"github.com/ashwinath/moneybags/pkg/db"
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
	baseDB, err := db.NewBaseDB(c.PostgresDb, logger)
	if err != nil {
		sugar.Fatalf("Failed to initialise Base DB, %v", err)
	}
	defer baseDB.Close()

	transactionDB, err := db.NewTransactionDB(baseDB)
	if err != nil {
		sugar.Fatalf("Failed to initialise Transaction DB, %v", err)
	}

	// Load framework
	fw := framework.New(c, sugar, map[string]any{
		db.TransactionDatabaseName: transactionDB,
	})

	// Load modules
	telegram, err := modules.NewTelegramModule(fw)
	if err != nil {
		sugar.Fatalf("Failed to initialise telegram module, %v", err)
	}

	// Run app
	app := framework.NewApp(sugar, telegram)

	ctx, cancel := context.WithCancel(context.Background())
	framework.ListenForSignal(cancel, sugar)

	app.Run(ctx)
}
