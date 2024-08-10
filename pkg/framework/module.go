package framework

import (
	"context"

	"go.uber.org/zap"
)

type Module interface {
	Start(ctx context.Context)
	Name() string
}

type App struct {
	logger  *zap.SugaredLogger
	modules map[string]Module
}

func NewApp(logger *zap.SugaredLogger, modules ...Module) *App {
	moduleMap := map[string]Module{}
	for _, m := range modules {
		moduleMap[m.Name()] = m
	}
	return &App{
		modules: moduleMap,
		logger:  logger,
	}
}

func (a *App) Run(ctx context.Context) {
	a.logger.Info("Starting modules")
	for _, m := range a.modules {
		a.logger.Info("Starting module: %s", m.Name())
		go m.Start(ctx)
	}
	<-ctx.Done()
	a.logger.Info("Stopping apps")
}
