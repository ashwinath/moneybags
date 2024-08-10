package financials

import (
	"github.com/ashwinath/moneybags/pkg/db"
	"github.com/ashwinath/moneybags/pkg/framework"
)

type StocksLoader struct {
	fw           framework.FW
	tradeDB      db.TradeDB
	symbolDB     db.SymbolDB
	alphavantage Alphavantage
}

func NewStocksLoader(fw framework.FW, alphavantage Alphavantage) Loader {
	return &StocksLoader{
		fw:           fw,
		tradeDB:      fw.GetDB(db.TradeDatabaseName).(db.TradeDB),
		symbolDB:     fw.GetDB(db.SymbolDatabaseName).(db.SymbolDB),
		alphavantage: alphavantage,
	}
}

func (l *StocksLoader) Load() error {
	if err := l.processSymbols(); err != nil {
		return err
	}
	return nil
}

func (l *StocksLoader) processSymbols() error {
	if err := l.processStockSymbols(); err != nil {
		return err
	}

	if err := l.processCurrencySymbols(); err != nil {
		return err
	}

	return nil
}

func (l *StocksLoader) processStockSymbols() error {
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

func (l *StocksLoader) processCurrencySymbols() error {
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
