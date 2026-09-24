package financials

import (
	"context"
	"fmt"
	"time"

	"github.com/ashwinath/simple/client"
)

const yahooBaseURL = "https://query1.finance.yahoo.com"

type yahoo struct{}

// NewYahoo creates a new Yahoo Finance market data provider.
func NewYahoo() MarketDataProvider {
	return &yahoo{}
}

func (y *yahoo) headers() map[string]string {
	return map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
	}
}

type yahooSearchResult struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"shortname"`
	Currency string `json:"currency"`
}

func (y *yahoo) GetSymbolInfo(symbol string) (*SymbolInfo, error) {
	url := fmt.Sprintf("%s/v8/finance/chart/%s?interval=1d&range=1d", yahooBaseURL, symbol)

	var res yahooChartResponse
	err := client.HTTPGet(context.TODO(), url, y.headers(), &res)
	if err != nil {
		return nil, fmt.Errorf("could not get symbol info (%s) from yahoo (%s): %s", symbol, url, err)
	}

	if len(res.Chart.Result) == 0 {
		return nil, fmt.Errorf("no results for symbol (%s) from yahoo", symbol)
	}

	result := res.Chart.Result[0]
	currency := "USD"
	if result.Meta.Currency != "" {
		currency = result.Meta.Currency
	}

	return &SymbolInfo{
		Symbol:   symbol,
		Currency: currency,
	}, nil
}

func (y *yahoo) GetCurrencyHistory(from string, to string) (map[string]OHLC, error) {
	ticker := fmt.Sprintf("%s%s=X", from, to)
	url := fmt.Sprintf("%s/v8/finance/chart/%s?interval=1d&range=5y", yahooBaseURL, ticker)

	var res yahooChartResponse
	err := client.HTTPGet(context.TODO(), url, y.headers(), &res)
	if err != nil {
		return nil, fmt.Errorf("could not get currency history (%s->%s) from yahoo (%s): %s", from, to, url, err)
	}

	if len(res.Chart.Result) == 0 {
		return nil, fmt.Errorf("currency history (%s->%s) result was empty from yahoo", from, to)
	}

	return convertYahooChartToOHLC(res.Chart.Result[0])
}

func (y *yahoo) GetStockHistory(symbol string) (map[string]OHLC, error) {
	url := fmt.Sprintf("%s/v8/finance/chart/%s?interval=1d&range=5y", yahooBaseURL, symbol)

	var res yahooChartResponse
	err := client.HTTPGet(context.TODO(), url, y.headers(), &res)
	if err != nil {
		return nil, fmt.Errorf("could not get stock history (%s) from yahoo (%s): %s", symbol, url, err)
	}

	if len(res.Chart.Result) == 0 {
		return nil, fmt.Errorf("stock history (%s) result was empty from yahoo", symbol)
	}

	return convertYahooChartToOHLC(res.Chart.Result[0])
}

type yahooChartResponse struct {
	Chart struct {
		Result []yahooChartResult `json:"result"`
		Error  *yahooChartError   `json:"error"`
	} `json:"chart"`
}

type yahooChartError struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type yahooChartResult struct {
	Meta struct {
		Symbol   string `json:"symbol"`
		Currency string `json:"currency"`
	} `json:"meta"`
	Timestamp  []int64 `json:"timestamp"`
	Indicators struct {
		Quote []struct {
			Open   []*float64 `json:"open"`
			High   []*float64 `json:"high"`
			Low    []*float64 `json:"low"`
			Close  []*float64 `json:"close"`
			Volume []*int64   `json:"volume"`
		} `json:"quote"`
	} `json:"indicators"`
}

func convertYahooChartToOHLC(result yahooChartResult) (map[string]OHLC, error) {
	if len(result.Timestamp) == 0 {
		return nil, fmt.Errorf("no timestamp data for symbol %s", result.Meta.Symbol)
	}

	if len(result.Indicators.Quote) == 0 {
		return nil, fmt.Errorf("no quote data for symbol %s", result.Meta.Symbol)
	}

	quote := result.Indicators.Quote[0]
	ohlcs := make(map[string]OHLC, len(result.Timestamp))

	for i, ts := range result.Timestamp {
		if i >= len(quote.Open) || i >= len(quote.High) || i >= len(quote.Low) || i >= len(quote.Close) {
			break
		}

		o := quote.Open[i]
		h := quote.High[i]
		l := quote.Low[i]
		c := quote.Close[i]

		if o == nil || h == nil || l == nil || c == nil {
			continue
		}

		date := time.Unix(ts, 0).UTC().Format(time.DateOnly)
		ohlcs[date] = OHLC{
			Open:  *o,
			High:  *h,
			Low:   *l,
			Close: *c,
		}
	}

	if len(ohlcs) == 0 {
		return nil, fmt.Errorf("no valid OHLC data for symbol %s", result.Meta.Symbol)
	}

	return ohlcs, nil
}
