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
	stockDB        db.StockDB
	portfolioDB    db.PortfolioDB
	alphavantage   Alphavantage
}

func NewStocksLoader(fw framework.FW, alphavantage Alphavantage) Loader {
	return &stocksLoader{
		fw:             fw,
		tradeDB:        fw.GetDB(db.TradeDatabaseName).(db.TradeDB),
		symbolDB:       fw.GetDB(db.SymbolDatabaseName).(db.SymbolDB),
		exchangeRateDB: fw.GetDB(db.ExchangeRateDatabaseName).(db.ExchangeRateDB),
		stockDB:        fw.GetDB(db.StockDatabaseName).(db.StockDB),
		portfolioDB:    fw.GetDB(db.PortfolioDatabaseName).(db.PortfolioDB),
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

	if err := l.processStocks(); err != nil {
		return fmt.Errorf("Failed to load process stocks: %s", err)
	}

	if err := l.calculatePortfolio(); err != nil {
		return fmt.Errorf("Failed to calculate portfolio: %s", err)
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

func (l *stocksLoader) processStocks() error {
	stockSymbols, err := l.symbolDB.GetStocks()
	if err != nil {
		return err
	}

	for _, stock := range stockSymbols {
		if err := l.processStock(stock); err != nil {
			return err
		}
	}

	return nil
}

func (l *stocksLoader) processStock(symbol db.Symbol) error {
	isCompact := symbol.LastProcessedDate != nil
	history, err := l.alphavantage.GetStockHistory(symbol.Symbol, isCompact)
	if err != nil {
		return err
	}

	var lastProcessedDate *time.Time
	stockHistory := []*db.Stock{}

	for date, value := range history {
		d, err := utils.SetDateFromString(date)
		if err != nil {
			return fmt.Errorf("could not parse date (%s): %s", date, err)
		}

		if lastProcessedDate == nil || lastProcessedDate.Before(d) {
			lastProcessedDate = &d
		}

		stock := db.Stock{
			TradeDate: d,
			Symbol:    symbol.Symbol,
			Price:     value.Close,
		}
		stockHistory = append(stockHistory, &stock)
	}

	if err := l.stockDB.BulkAdd(stockHistory); err != nil {
		return err
	}

	if err := l.symbolDB.UpdateLastProcessedDate(symbol.Symbol, *lastProcessedDate); err != nil {
		return err
	}

	return nil
}

func (l *stocksLoader) calculatePortfolio() error {
	symbols, err := l.symbolDB.GetStocks()
	if err != nil {
		return fmt.Errorf("failed to retrieve symbols from db.")
	}

	for _, symbol := range symbols {
		trades, err := l.tradeDB.GetTradesSorted(symbol.Symbol)
		if err != nil {
			return fmt.Errorf("failed to retrieve trades with symbol (%s): %s", symbol.Symbol, err)
		}

		partialPortfolios := []db.Portfolio{}

		// First pass to fill active trading parts
		// NAV and SimpleReturns are filled later.
		for _, trade := range trades {
			exchangeRate, err := l.getCurrencyRate(trade.DatePurchased.Time, *symbol.BaseCurrency)
			if err != nil {
				return err
			}

			tradeMultiplier := 1.0
			if trade.TradeType != "buy" {
				tradeMultiplier = -1.0
			}

			portfolio := db.Portfolio{
				TradeDate: trade.DatePurchased.Time,
				Symbol:    symbol.Symbol,
			}
			if len(partialPortfolios) == 0 {
				portfolio.Principal = trade.PriceEach * trade.Quantity * exchangeRate
				portfolio.Quantity = trade.Quantity
			} else {
				lastPortfolio := partialPortfolios[len(partialPortfolios)-1]
				portfolio.Principal = lastPortfolio.Principal + (trade.PriceEach * trade.Quantity * exchangeRate * tradeMultiplier)
				portfolio.Quantity = lastPortfolio.Quantity + (trade.Quantity * tradeMultiplier)
			}

			partialPortfolios = append(partialPortfolios, portfolio)
		}

		// Second pass to update all gaps in non active trading days
		portfolioMap := map[time.Time]db.Portfolio{}
		for _, partial := range partialPortfolios {
			// There might be multiple trades in a single day for each symbol, we need to combine them
			if portfolio, ok := portfolioMap[partial.TradeDate]; ok {
				// Override with latest value
				partial.Quantity = portfolio.Quantity
				partial.Principal = portfolio.Principal
			}
			portfolioMap[partial.TradeDate] = partial
		}

		allPortfolios := []db.Portfolio{}
		currentDate := partialPortfolios[0].TradeDate
		tomorrow := time.Now().AddDate(0, 0, 1)

		for currentDate.Before(tomorrow) {
			exchangeRate, err := l.getCurrencyRate(currentDate, *symbol.BaseCurrency)
			if err != nil {
				return err
			}

			price, err := l.getStockPrice(currentDate, symbol.Symbol)
			if err != nil {
				return err
			}

			var previousPortfolio db.Portfolio
			if p, ok := portfolioMap[currentDate]; ok {
				previousPortfolio = p
			} else {
				// guaranteed to have an element
				previousPortfolio = allPortfolios[len(allPortfolios)-1]
			}

			newPortfolio := db.Portfolio{
				TradeDate: currentDate,
				Symbol:    symbol.Symbol,
				Principal: previousPortfolio.Principal,
				Quantity:  previousPortfolio.Quantity,
				NAV:       previousPortfolio.Quantity * price * exchangeRate,
			}
			newPortfolio.SimpleReturns = (newPortfolio.NAV - newPortfolio.Principal) / newPortfolio.Principal
			allPortfolios = append(allPortfolios, newPortfolio)
			currentDate = currentDate.AddDate(0, 0, 1)
		}

		if err := l.portfolioDB.BulkAdd(allPortfolios); err != nil {
			return fmt.Errorf("failed to bulk insert into porfolio: %s", err)
		}
	}

	return nil
}

func (l *stocksLoader) getCurrencyRate(date time.Time, currency string) (float64, error) {
	for {
		if val, err := l.exchangeRateDB.GetExchangeRateByDate(date, currency); err == nil {
			return val, nil
		}
		date = date.AddDate(0, 0, -1)
	}
}

func (l *stocksLoader) getStockPrice(date time.Time, symbol string) (float64, error) {
	for {
		if val, err := l.stockDB.GetStockPrice(date, symbol); err == nil {
			return val, nil
		}

		date = date.AddDate(0, 0, -1)
	}
}
