package financials

// Loader defines the interface for data loaders.
type Loader interface {
	Load() error
	Name() string
}

// MarketDataProvider defines the interface for fetching market data
// (symbols, currency history, stock history) from external providers.
type MarketDataProvider interface {
	GetSymbolInfo(symbol string) (*SymbolInfo, error)
	GetCurrencyHistory(from string, to string) (map[string]OHLC, error)
	GetStockHistory(symbol string) (map[string]OHLC, error)
}

// SymbolInfo represents information about a trading symbol.
type SymbolInfo struct {
	Symbol   string
	Currency string
}

// OHLC represents Open, High, Low, Close price data for a given date.
type OHLC struct {
	Open  float64
	High  float64
	Low   float64
	Close float64
}
