package modules

import (
	"context"

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
}
