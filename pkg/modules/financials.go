package modules

import (
	"context"
	"fmt"
	"time"

	"github.com/ashwinath/moneybags/pkg/financials"
	"github.com/ashwinath/moneybags/pkg/framework"
)

type FinancialsModule struct {
	fw      framework.FW
	loaders []financials.Loader
}

func NewFinancialsModule(fw framework.FW, alphavantage financials.Alphavantage) (framework.Module, error) {
	if fw.GetConfig().FinancialsConfig.AlphavantageApiKey == "" {
		return nil, fmt.Errorf("alphavantageApiKey not set")
	}
	return &FinancialsModule{
		fw: fw,
		loaders: []financials.Loader{
			financials.NewCSVLoader(fw),
			financials.NewTransactionLoader(fw),
			financials.NewStocksLoader(fw, alphavantage),
			financials.NewInvestmentsLoader(fw),
			financials.NewSharedExpenseLoader(fw),
			financials.NewAverageExpenditureLoader(fw),
			financials.NewMortgageLoader(fw),
			financials.NewHouseAssetLoader(fw),
			financials.NewCarLoader(fw),
		},
	}, nil
}

func (m *FinancialsModule) Name() string {
	return "financials"
}

func (m *FinancialsModule) Start(ctx context.Context) {
	framework.RunInterval(
		ctx,
		time.Duration(m.fw.GetConfig().FinancialsConfig.RunIntervalInHours)*time.Hour,
		m.run,
	)
}

func (m *FinancialsModule) run() {
	hasError := false
	m.fw.GetLogger().Infof("running one round of financials module")
	m.fw.TimeFunction("Financials Module", func() {
		for _, loader := range m.loaders {
			m.fw.TimeFunction(fmt.Sprintf("[%s] loader", loader.Name()), func() {
				if err := loader.Load(); err != nil {
					m.fw.GetLogger().Errorf("Failed to run loader: %s", err)
					hasError = true
					return
				}
			})
			if hasError {
				break
			}
		}
	})
	m.fw.GetLogger().Infof("finished running one round of financials module")
}
