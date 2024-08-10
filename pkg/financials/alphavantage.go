package financials

import (
	"fmt"

	"github.com/ashwinath/moneybags/pkg/framework"
)

type Alphavantage interface {
	GetSymbolFromAlphavantage(symbol string) (*AlphavantageSymbol, error)
}

type alphavantage struct {
	apiKey string
}

func NewAlphavantage(apiKey string) Alphavantage {
	return &alphavantage{
		apiKey: apiKey,
	}
}

type symbolResult struct {
	BestMatches []AlphavantageSymbol `json:"bestMatches"`
}

type AlphavantageSymbol struct {
	Symbol   string `json:"1. symbol"`
	Currency string `json:"8. currency"`
}

func (a *alphavantage) GetSymbolFromAlphavantage(symbol string) (*AlphavantageSymbol, error) {
	url := fmt.Sprintf(
		"https://www.alphavantage.co/query?function=SYMBOL_SEARCH&keywords=%s&apikey=%s",
		symbol, a.apiKey,
	)
	res := symbolResult{}
	err := framework.RetrySimple(func() error {
		return framework.HTTPGet(url, &res)
	})
	if err != nil {
		return nil, fmt.Errorf("Could not get symbol (%s) result from alphavantage: %s", symbol, err)
	}

	// For testing we use demo key
	if symbol != "tesco" {
		if len(res.BestMatches) != 1 {
			return nil, fmt.Errorf("Could not get symbol (%s) result from alphavantage, there was not equal to 1 results, length = %d", symbol, len(res.BestMatches))
		}
	}

	sym := res.BestMatches[0]
	return &sym, nil
}
