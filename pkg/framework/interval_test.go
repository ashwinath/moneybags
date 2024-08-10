package framework

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunInterval(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	counter := 0
	var m sync.Mutex
	counterFunction := func() {
		m.Lock()
		counter++
		m.Unlock()
	}

	go RunInterval(ctx, 100*time.Millisecond, counterFunction)
	time.Sleep(300 * time.Millisecond)
	cancel()

	m.Lock()
	assert.Greater(t, counter, 1)
	m.Unlock()
}
