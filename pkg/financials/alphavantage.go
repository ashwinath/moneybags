package financials

import (
	"fmt"
	"strconv"

	"github.com/ashwinath/moneybags/pkg/framework"
)

type Alphavantage interface {
	GetSymbolFromAlphavantage(symbol string) (*AlphavantageSymbol, error)
	GetCurrencyHistory(from string, to string, isCompact bool) (map[string]OHLC, error)
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

type fxDailyResult struct {
	Results map[string]alphavantageOHLC `json:"Time Series FX (Daily)"`
}

type alphavantageOHLC struct {
	Open  string `json:"1. open"`
	High  string `json:"2. high"`
	Low   string `json:"3. low"`
	Close string `json:"4. close"`
}

type OHLC struct {
	Open  float64
	High  float64
	Low   float64
	Close float64
}

func (a *alphavantage) GetCurrencyHistory(from string, to string, isCompact bool) (map[string]OHLC, error) {
	outputSize := "full"
	if isCompact {
		outputSize = "compact"
	}
	url := fmt.Sprintf(
		"https://www.alphavantage.co/query?function=FX_DAILY&from_symbol=%s&to_symbol=%s&outputsize=%s&apikey=%s",
		from, to, outputSize, a.apiKey,
	)

	fmt.Println(url)

	res := fxDailyResult{}
	err := framework.RetrySimple(func() error {
		return framework.HTTPGet(url, &res)
	})
	if err != nil {
		return nil, fmt.Errorf("Could not get currency history (%s->%s) result from alphavantage: %s", from, to, err)
	}

	returnMap := map[string]OHLC{}
	for key, value := range res.Results {
		o, err := strconv.ParseFloat(value.Open, 64)
		if err != nil {
			return nil, fmt.Errorf(
				"could not convert exchange rate (%s -> %s) open to float64: %s",
				from, to, err,
			)
		}
		h, err := strconv.ParseFloat(value.High, 64)
		if err != nil {
			return nil, fmt.Errorf(
				"could not convert exchange rate (%s -> %s) high to float64: %s",
				from, to, err,
			)
		}
		l, err := strconv.ParseFloat(value.Low, 64)
		if err != nil {
			return nil, fmt.Errorf(
				"could not convert exchange rate (%s -> %s) low to float64: %s",
				from, to, err,
			)
		}
		c, err := strconv.ParseFloat(value.Close, 64)
		if err != nil {
			return nil, fmt.Errorf(
				"could not convert exchange rate (%s -> %s) close to float64: %s",
				from, to, err,
			)
		}
		returnMap[key] = OHLC{
			Open:  o,
			High:  h,
			Low:   l,
			Close: c,
		}
	}

	return returnMap, nil
}
