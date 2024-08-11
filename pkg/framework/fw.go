package framework

import (
	"time"

	"github.com/ashwinath/moneybags/pbgo/configpb"
	"go.uber.org/zap"
)

type FW interface {
	GetConfig() *configpb.Config
	GetDB(name string) any
	GetLogger() *zap.SugaredLogger
	TimeFunction(functionName string, fn func())
}

type Framework struct {
	config    *configpb.Config
	databases map[string]any
	logger    *zap.SugaredLogger
}

func New(config *configpb.Config, logger *zap.SugaredLogger, dbs map[string]any) FW {
	logger.Info("Initialising Framework")
	return &Framework{
		config:    config,
		databases: dbs,
		logger:    logger,
	}
}

func (fw *Framework) GetConfig() *configpb.Config {
	return fw.config
}

func (fw *Framework) GetLogger() *zap.SugaredLogger {
	return fw.logger
}

func (fw *Framework) GetDB(name string) any {
	return fw.databases[name]
}

func (fw *Framework) TimeFunction(functionName string, fn func()) {
	start := time.Now()
	fn()
	fw.logger.Infof("%s ran for %s", functionName, time.Since(start))
}
