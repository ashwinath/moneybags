package main

import (
	"flag"

	"github.com/ashwinath/moneybags/pkg/config"
	"github.com/ashwinath/moneybags/pkg/db"
	"github.com/ashwinath/moneybags/pkg/framework"
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
	_ = framework.New(c, sugar, map[string]any{
		"transaction": transactionDB,
	})

	// TODO: load modules
}
