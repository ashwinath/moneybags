package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	usd = "USD"
)

func TestCheckIfSymbolExists(t *testing.T) {
	err := RunTest(func(db *DB) {
		symbolDB, err := NewSymbolDB(db)
		assert.Nil(t, err)

		symbol := Symbol{
			SymbolType:   SymbolTypeStock,
			Symbol:       "world",
			BaseCurrency: &usd,
		}

		err = symbolDB.Insert(&symbol)
		assert.Nil(t, err)

		exists, err := symbolDB.CheckIfSymbolExists("world")
		assert.Nil(t, err)
		assert.True(t, exists)

		exists, err = symbolDB.CheckIfSymbolExists("non-existent")
		assert.Nil(t, err)
		assert.False(t, exists)
	})
	assert.Nil(t, err)
}

func TestGetDistinctCurrencies(t *testing.T) {
	err := RunTest(func(db *DB) {
		symbolDB, err := NewSymbolDB(db)
		assert.Nil(t, err)

		symbols := []Symbol{
			{
				SymbolType:   SymbolTypeStock,
				Symbol:       "world",
				BaseCurrency: &usd,
			},
			{
				SymbolType:   SymbolTypeStock,
				Symbol:       "hello",
				BaseCurrency: &usd,
			},
			{
				SymbolType: SymbolTypeCurrency,
				Symbol:     "usd",
			},
		}

		for _, symbol := range symbols {
			err = symbolDB.Insert(&symbol)
			assert.Nil(t, err)
		}

		currencies, err := symbolDB.GetDistinctCurrencies()
		assert.Nil(t, err)
		assert.Len(t, currencies, 1)
	})
	assert.Nil(t, err)
}
