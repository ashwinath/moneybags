package framework

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

type testModule struct{}

func (m *testModule) Start(ctx context.Context) {
	<-ctx.Done()
}

func (m *testModule) Name() string {
	return "test-module"
}

func newTestModule() Module {
	return &testModule{}
}

func TestFrameworkApp(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	a := NewApp(sugar, newTestModule())
	ctx, cancel := context.WithCancel(context.Background())
	finished := make(chan struct{}, 1)
	go func() {
		a.Run(ctx)
		finished <- struct{}{}
	}()

	cancel()
	<-finished
}
