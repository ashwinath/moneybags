package financials

import (
	"testing"

	database "github.com/ashwinath/moneybags/pkg/db"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	err := database.RunTest(func(db *database.DB) {
		fw := createFW(t, db)
		loader := NewCSVLoader(fw)
		err := loader.Load()
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
