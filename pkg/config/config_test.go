package config

import (
	"path"
	"testing"

	"github.com/ashwinath/moneybags/pbgo/configpb"
	"github.com/ashwinath/moneybags/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	p := path.Join(utils.GetLocalRepoLocation(), "./pkg/config/testdata/config.yaml")

	c, err := New(p)
	assert.Nil(t, err)

	assert.Equal(t, &configpb.Config{
		PostgresDb: &configpb.PostgresDB{
			Host:     "127.0.0.1",
			User:     "postgres",
			Password: "very_secure",
			DbName:   "postgres",
			Port:     5432,
		},
		TelegramConfig: &configpb.TelegramConfig{
			ApiKey:      "very secret",
			Debug:       true,
			AllowedUser: "hello",
		},
		FinancialsData: &configpb.FinancialsData{
			AssetsCsvFilepath:         "sample/assets.csv",
			ExpensesCsvFilepath:       "sample/expenses.csv",
			IncomeCsvFilepath:         "sample/income.csv",
			SharedExpensesCsvFilepath: "sample/shared_expenses.csv",
			TradesCsvFilepath:         "sample/trades.csv",
		},
	}, c)
}
