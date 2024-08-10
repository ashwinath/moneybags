package modules

import (
	"context"
	"time"

	"github.com/ashwinath/moneybags/pkg/framework"
)

type FinancialsModule struct {
	fw framework.FW
}

func NewFinancialsModule(fw framework.FW) (framework.Module, error) {
	return &FinancialsModule{
		fw: fw,
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

func (m *FinancialsModule) run() {}
