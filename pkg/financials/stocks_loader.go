package financials

import (
	"fmt"
	"time"

	"github.com/ashwinath/moneybags/pkg/db"
	"github.com/ashwinath/moneybags/pkg/framework"
	"github.com/ashwinath/moneybags/pkg/utils"
)

type stocksLoader struct {
	fw             framework.FW
	tradeDB        db.TradeDB
	symbolDB       db.SymbolDB
	exchangeRateDB db.ExchangeRateDB
	alphavantage   Alphavantage
}

func NewStocksLoader(fw framework.FW, alphavantage Alphavantage) Loader {
	return &stocksLoader{
		fw:             fw,
		tradeDB:        fw.GetDB(db.TradeDatabaseName).(db.TradeDB),
		symbolDB:       fw.GetDB(db.SymbolDatabaseName).(db.SymbolDB),
		exchangeRateDB: fw.GetDB(db.ExchangeRateDatabaseName).(db.ExchangeRateDB),
		alphavantage:   alphavantage,
	}
}

func (l *stocksLoader) Load() error {
	if err := l.processSymbols(); err != nil {
		return fmt.Errorf("Failed to load process symbols: %s", err)
	}

	if err := l.processCurrencies(); err != nil {
		return fmt.Errorf("Failed to load process currencies: %s", err)
	}

	return nil
}

func (l *stocksLoader) processSymbols() error {
	if err := l.processStockSymbols(); err != nil {
		return err
	}

	if err := l.processCurrencySymbols(); err != nil {
		return err
	}

	return nil
}

func (l *stocksLoader) processStockSymbols() error {
	tradeSymbols, err := l.tradeDB.GetUniqueSymbols()
	if err != nil {
		return nil
	}

	for _, s := range tradeSymbols {
		exists, err := l.symbolDB.CheckIfSymbolExists(s)
		if err != nil {
			return err
		}

		if !exists {
			aSym, err := l.alphavantage.GetSymbolFromAlphavantage(s)
			if err != nil {
				return err
			}
			sym := db.Symbol{
				SymbolType:   db.SymbolTypeStock,
				Symbol:       s,
				BaseCurrency: &aSym.Currency,
			}
			if err := l.symbolDB.Insert(&sym); err != nil {
				return err
			}
		}

	}

	return nil
}

func (l *stocksLoader) processCurrencySymbols() error {
	currencies, err := l.symbolDB.GetDistinctCurrencies()
	if err != nil {
		return err
	}

	for _, currency := range currencies {
		exists, err := l.symbolDB.CheckIfSymbolExists(currency)
		if err != nil {
			return err
		}

		if !exists {
			sym := db.Symbol{
				SymbolType: db.SymbolTypeCurrency,
				Symbol:     currency,
			}
			if err := l.symbolDB.Insert(&sym); err != nil {
				return err
			}
		}
	}
	return nil
}

func (l *stocksLoader) processCurrencies() error {
	currencySymbols, err := l.symbolDB.GetCurrencies()
	if err != nil {
		return err
	}

	for _, currency := range currencySymbols {
		if err := l.processCurrency(currency); err != nil {
			return err
		}
	}

	return nil
}

func (l *stocksLoader) processCurrency(symbol db.Symbol) error {
	isCompact := symbol.LastProcessedDate != nil
	history, err := l.alphavantage.GetCurrencyHistory(symbol.Symbol, "SGD", isCompact)
	if err != nil {
		return err
	}

	var lastProcessedDate *time.Time
	currencyHistory := []*db.ExchangeRate{}

	for date, value := range history {
		d, err := utils.SetDateFromString(date)
		if err != nil {
			return fmt.Errorf("could not parse date (%s): %s", date, err)
		}

		if lastProcessedDate == nil || lastProcessedDate.Before(d) {
			lastProcessedDate = &d
		}

		er := db.ExchangeRate{
			TradeDate: d,
			Symbol:    symbol.Symbol,
			Price:     value.Close,
		}
		currencyHistory = append(currencyHistory, &er)
	}

	if err := l.exchangeRateDB.BulkAdd(currencyHistory); err != nil {
		return err
	}

	if err := l.symbolDB.UpdateLastProcessedDate(symbol.Symbol, *lastProcessedDate); err != nil {
		return err
	}

	return nil
}
