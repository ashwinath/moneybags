package framework

import (
	"context"
	"time"
)

func RunInterval(ctx context.Context, t time.Duration, fn func()) {
	sleepDone := make(chan struct{})
	for {
		go func() {
			fn()
			time.Sleep(t)
			sleepDone <- struct{}{}
		}()
		select {
		case <-ctx.Done():
			break
		case <-sleepDone:
			continue
		}
	}
}
