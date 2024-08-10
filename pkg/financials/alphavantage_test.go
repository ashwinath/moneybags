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
