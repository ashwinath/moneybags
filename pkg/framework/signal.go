package framework

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func ListenForSignal(cancel context.CancelFunc, logger *zap.SugaredLogger) {
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-signalChannel
		switch sig {
		case os.Interrupt:
			logger.Infof("received interrupt signal")
			cancel()
		case syscall.SIGTERM:
			logger.Infof("received sigterm")
			cancel()
		}
	}()
}
