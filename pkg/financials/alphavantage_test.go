package financials

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAlphavantageSymbol(t *testing.T) {
	av := NewAlphavantage("demo")

	sym, err := av.GetSymbolFromAlphavantage("tesco")
	assert.Nil(t, err)
	assert.Equal(t, "TSCO.LON", sym.Symbol)
	assert.Equal(t, "GBX", sym.Currency)
}

func TestGetCurrencyHistory(t *testing.T) {
	av := NewAlphavantage("demo")

	ohlcs, err := av.GetCurrencyHistory("EUR", "USD", false)
	assert.Nil(t, err)

	assert.Greater(t, len(ohlcs), 1)

	// Test one value
	for _, value := range ohlcs {
		assert.Greater(t, value.Open, 0.0000001)
		assert.Greater(t, value.High, 0.0000001)
		assert.Greater(t, value.Low, 0.0000001)
		assert.Greater(t, value.Close, 0.0000001)
		break
	}
}
