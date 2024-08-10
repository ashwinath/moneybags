package modules

import (
	"context"
	"time"

	"github.com/ashwinath/moneybags/pkg/financials"
	"github.com/ashwinath/moneybags/pkg/framework"
)

type FinancialsModule struct {
	fw      framework.FW
	loaders []financials.Loader
}

func NewFinancialsModule(fw framework.FW, alphavantage financials.Alphavantage) (framework.Module, error) {
	return &FinancialsModule{
		fw: fw,
		loaders: []financials.Loader{
			financials.NewCSVLoader(fw),
			financials.NewTelegramLoader(fw),
			financials.NewStocksLoader(fw, alphavantage),
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
	for _, loader := range m.loaders {
		if err := loader.Load(); err != nil {
			m.fw.GetLogger().Errorf("Failed to run loader: %s", err)
		}
	}
}
